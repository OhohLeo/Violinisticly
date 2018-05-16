package client

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"unsafe"

	"github.com/ohohleo/go-dsp/fft"
	"github.com/xthexder/go-jack"
)

// Use global ports to avoid CGO issue
var portsIn []*jack.Port
var mutexIn *sync.Mutex
var onInput func([]jack.AudioSample)

var portsOut []*jack.Port
var mutexOut *sync.Mutex
var onOutput func([]jack.AudioSample)

type Client struct {
	client *jack.Client
	name   string

	portsIn  []*jack.Port
	portsOut []*jack.Port

	links map[string]string

	inChan       chan []float32
	inFFTChan    chan []complex128
	inFFTBuffer  []float64
	outChan      chan []float32
	outFFTChan   chan []complex128
	outFFTBuffer []float64

	stopChan chan struct{}
}

func New(name string, links map[string]string) (c *Client, err error) {

	if name == "" {
		return nil, fmt.Errorf("name required")
	}

	if links == nil {
		return nil, fmt.Errorf("links required")
	}

	c = &Client{
		links: links,
	}

	var status int
	c.client, status = jack.ClientOpen(name, jack.NoStartServer)
	if err = jack.Strerror(status); err != nil {
		err = fmt.Errorf("client open issue %s", err.Error())
		return
	}

	status = c.client.SetProcessCallback(onProcess)
	if err = jack.Strerror(status); err != nil {
		err = fmt.Errorf("process callback issue %s", err.Error())
		return
	}

	c.client.OnShutdown(shutdown)

	c.name = name

	return
}

func (c *Client) Start() (err error) {

	// Activate client
	status := c.client.Activate()
	if err = jack.Strerror(status); err != nil {
		err = fmt.Errorf("activate issue %s", err.Error())
		return
	}

	// Check realtime
	if c.client.IsRealtime() {
		log.Println("realtime OK")
	}

	name := c.name + ":"
	nameSize := len(name)

	for src, dst := range c.links {

		// Register input ports
		if strings.HasPrefix(dst, name) {

			log.Println("register input:" + dst[nameSize:])

			portIn := c.client.PortRegister(
				dst[nameSize:], jack.DEFAULT_AUDIO_TYPE, jack.PortIsInput, 0)
			c.portsIn = append(c.portsIn, portIn)
		}

		// Register output ports
		if strings.HasPrefix(src, name) {

			log.Println("register output:" + dst[nameSize:])

			portOut := c.client.PortRegister(
				src[nameSize:], jack.DEFAULT_AUDIO_TYPE, jack.PortIsOutput, 0)
			c.portsOut = append(c.portsOut, portOut)
		}

		// Establish connections
		status := c.client.Connect(src, dst)
		if err = jack.Strerror(status); err != nil {
			log.Printf("link %s => %s KO", src, dst)
			err = fmt.Errorf("connect issue %s", err.Error())
			return
		}

		log.Printf("link %s => %s OK", src, dst)
	}

	// Set reference into global
	portsIn = c.portsIn
	portsOut = c.portsOut

	log.Println("START " + c.name)

	// Wait until stop
	c.stopChan = make(chan struct{})
	<-c.stopChan

	return
}

func (c *Client) Stop() (err error) {

	log.Println("STOP " + c.name)

	// Stop channels
	if c.inChan != nil {

		mutexIn.Lock()
		onInput = nil
		mutexIn.Unlock()

		close(c.inChan)
	}

	if c.outChan != nil {

		mutexOut.Lock()
		onOutput = nil
		mutexOut.Unlock()

		close(c.outChan)
	}

	if c.stopChan != nil {
		close(c.stopChan)
	}

	// Establish disconnect
	for src, dst := range c.links {
		status := c.client.Disconnect(src, dst)
		if err = jack.Strerror(status); err != nil {
			err = fmt.Errorf("disconnect issue %s", err.Error())
			return
		}
	}

	return jack.Strerror(c.client.Close())
}

func (c *Client) GetInput() (chan []float32, chan []complex128) {

	c.inChan = make(chan []float32)
	c.inFFTChan = make(chan []complex128)
	mutexIn = new(sync.Mutex)
	onInput = c.onInput

	return c.inChan, c.inFFTChan
}

func (c *Client) onInput(samples []jack.AudioSample) {

	// Convert []AudioSample => []float32
	values := *(*[]float32)(unsafe.Pointer(&samples))
	go handleFFT(values, c.inFFTChan)
	// Send raw values
	select {
	case c.inChan <- values:
		return
	}
}

func (c *Client) GetOutput() (chan []float32, chan []complex128) {

	c.outChan = make(chan []float32)
	c.outFFTChan = make(chan []complex128)
	mutexOut = new(sync.Mutex)
	onOutput = c.onOutput

	return c.outChan, c.outFFTChan
}

func (c *Client) onOutput(samples []jack.AudioSample) {

	// Convert []AudioSample => []float32
	values := *(*[]float32)(unsafe.Pointer(&samples))
	go handleFFT(values, c.outFFTChan)
	// Send raw values
	select {
	case c.outChan <- values:
		return
	}
}

func handleFFT(values []float32, fftChan chan []complex128) {

	// Convert into 64 bits
	inputs := make([]float64, len(values))
	for idx, value := range values {
		inputs[idx] = float64(value)
	}

	// Calculate FFT
	select {
	case fftChan <- fft.FFTReal(inputs):
		return
	}
}

func shutdown() {
	fmt.Println("Shutting down")
}

func onProcess(framesNb uint32) int {

	for portIdx, in := range portsIn {

		// Get samples input
		samplesIn := in.GetBuffer(framesNb)

		if onInput != nil {
			mutexIn.Lock()
			onInput(samplesIn)
			mutexIn.Unlock()
		}

		// Get samples output
		samplesOut := portsOut[portIdx].GetBuffer(framesNb)

		// Copy input into output
		for idx, sample := range samplesIn {
			samplesOut[idx] = sample
		}

		if onOutput != nil {
			mutexOut.Lock()
			onOutput(samplesOut)
			mutexOut.Unlock()
		}
	}

	return 0 // no error
}
