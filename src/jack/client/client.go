package client

import (
	"fmt"
	"log"
	"strings"

	"github.com/xthexder/go-jack"
)

// Use global ports to avoid CGO issue
var portsIn []*jack.Port
var portsOut []*jack.Port

type Client struct {
	client *jack.Client
	name   string

	portsIn  []*jack.Port
	portsOut []*jack.Port

	links map[string]string

	channel chan struct{}
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
	c.channel = make(chan struct{})
	<-c.channel

	return
}

func (c *Client) Stop() (err error) {

	log.Println("STOP " + c.name)

	// Stop channel
	if c.channel != nil {
		close(c.channel)
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

func shutdown() {
	fmt.Println("Shutting down")
}

func onProcess(framesNb uint32) int {

	for portIdx, in := range portsIn {

		// Get samples input
		samplesIn := in.GetBuffer(framesNb)

		// Get samples output
		samplesOut := portsOut[portIdx].GetBuffer(framesNb)

		// Copy input into output
		for idx, sample := range samplesIn {
			samplesOut[idx] = sample
		}
	}

	return 0 // no error
}
