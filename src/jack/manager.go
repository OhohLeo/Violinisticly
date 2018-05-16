package main

import (
	"fmt"
	"log"
	"math/cmplx"
	"os"
	"os/signal"

	"engo.io/engo"
	"github.com/ohohleo/violin/jack/client"
	"github.com/ohohleo/violin/jack/graphs"
)

const (
	DEFAULT_WIDTH  = 1000
	DEFAULT_HEIGHT = 800
)

type SoundManager struct {
	client  *client.Client
	signal  chan os.Signal
	display bool
}

func NewSoundManager(display bool) (s *SoundManager) {

	s = &SoundManager{
		signal:  make(chan os.Signal, 1),
		display: display,
	}

	signal.Notify(s.signal, os.Interrupt)

	return
}

func (s *SoundManager) Start(name string, links map[string]string) (err error) {

	s.client, err = client.New(name, links)
	if err != nil {
		return
	}

	go func() {
		signal := <-s.signal
		log.Println("Got signal:", signal)

		s.Stop()
	}()

	if s.display {

		go func() {
			if err = s.client.Start(); err != nil {
				panic(err)
			}
		}()

		// Can't use display inside go routine
		return s.StartDisplay()
	}

	return s.client.Start()
}

func (s *SoundManager) Stop() (err error) {

	if s.display {
		s.StopDisplay()
	}

	return s.client.Stop()
}

func (s *SoundManager) StartDisplay() error {

	// Load font
	fontPath := "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf"
	if err := engo.Files.Load(fontPath); err != nil {
		return err
	}

	input, inputFFT := s.client.GetInput()
	// output, outputFFT := s.client.GetOutput()

	inDisplayFFT := make(chan []float32)
	go handleFFT(inputFFT, inDisplayFFT, false)

	// outDisplayFFT := make(chan []float32)
	//go handleFFT(outputFFT, outDisplayFFT, false)

	scene, err := graphs.NewScene(
		fontPath,
		map[string]chan []float32{
			"input": input,
			//"output": output,
		},
		map[string]chan []float32{
			"input": inDisplayFFT,
			//"output": outDisplayFFT,
		},
		DEFAULT_WIDTH, DEFAULT_HEIGHT)
	if err != nil {
		return err
	}

	engo.Run(engo.RunOptions{
		Title:  "Graph",
		Width:  DEFAULT_WIDTH,
		Height: DEFAULT_HEIGHT,
	}, scene)

	return nil
}

func (s *SoundManager) StopDisplay() {
	engo.Exit()
}

func handleFFT(inFFT chan []complex128, outFFT chan []float32, debug bool) {

	var previousPhase float64

	for {
		fft, ok := <-inFFT
		if ok == false {
			close(outFFT)
			return
		}

		fftNb := len(fft)
		idx := fftNb / 5 // Final frequency 8785 Hz at 204
		//fmt.Printf("Final frequency %d %.3f\n", idx, 44100*float64(idx)/float64(fftNb))

		// Prepare values
		values := make([]float32, idx)

		var maxIdx int
		var maxValue, maxPhase, deltaPhase float64

		for idx > 0 {
			idx--
			value, phase := cmplx.Polar(fft[idx])

			if debug && value > maxValue {
				maxValue = float64(value)
				maxPhase = phase
				deltaPhase = phase - previousPhase
				maxIdx = idx
			}

			values[idx] = -float32(value) * 10
		}

		if debug {

			deltaFreq := deltaPhase * 44100 / float64(fftNb) * 0.314 / 2
			freq := 44100*float64(maxIdx)/float64(fftNb) + deltaFreq

			fmt.Printf("MAX at %04d nb:%d f:%.3f Hz delta f: %.3f value:%.3f delta phase:%.3f\n",
				maxIdx, fftNb, freq, deltaFreq, maxValue, deltaPhase)
			previousPhase = maxPhase
		}

		outFFT <- values
	}
}
