package main

import (
	"github.com/ohohleo/violin/api"
	"github.com/ohohleo/violin/input"
	"github.com/ohohleo/violin/opengl"
	"log"
	"net/http"
)

func main() {

	accelerometer, err := input.AccelGyroSerial("/dev/ttyACM0", 38400, false)
	if err != nil {
		log.Fatal(err)
	}

	err = api.New()

	if err != nil {
		log.Fatal(err)
	}

	err = api.NewStream(accelerometer)

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := opengl.CreateWindow(); err != nil {
			panic(err)
		}
	}()

	log.Println("Listening :5000 ...")
	http.ListenAndServe(":5000", nil)
}
