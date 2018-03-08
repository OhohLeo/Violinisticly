package main

import (
	"github.com/ohohleo/violin/jack/client"
	"log"
	"os"
	"os/signal"
)

type SoundMgmt struct {
	client *client.Client
	signal chan os.Signal
}

func NewSoundMgmt() (s *SoundMgmt) {

	s = &SoundMgmt{
		signal: make(chan os.Signal, 1),
	}

	signal.Notify(s.signal, os.Interrupt)

	return
}

func (s *SoundMgmt) Start(name string, links map[string]string) (err error) {

	s.client, err = client.New(name, links)
	if err != nil {
		return
	}

	go func() {
		signal := <-s.signal
		log.Println("Got signal:", signal)

		s.Stop()
	}()

	return s.client.Start()
}

func (s *SoundMgmt) Stop() (err error) {
	return s.client.Stop()
}

func main() {

	s := NewSoundMgmt()

	links := map[string]string{
		"PulseAudio JACK Sink:front-left":  "test:in_0",
		"PulseAudio JACK Sink:front-right": "test:in_1",
		"test:out_0":                       "system:playback_1",
		"test:out_1":                       "system:playback_2",
	}

	if err := s.Start("test", links); err != nil {
		panic(err)
	}

}
