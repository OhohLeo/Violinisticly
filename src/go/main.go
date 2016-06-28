package main

import (
	"log"
	"net/http"
	"ohohleo/accelerometer/api"
	"ohohleo/accelerometer/input"
)

func main() {

	accelerometer, err := input.AccelerometerSerial("/dev/ttyACM0", 115200)
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

	log.Println("HTTP Listening :5000 ...")
	http.ListenAndServe(":5000", nil)
}
