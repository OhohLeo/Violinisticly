package main

import (
	"flag"
)

func main() {

	display := flag.Bool("display", false, "activate display")
	flag.Parse()

	s := NewSoundManager(*display)

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
