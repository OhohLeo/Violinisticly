package main

import (
	// "github.com/ohohleo/violin/api"
	"fmt"
	"github.com/ohohleo/violin/input"
	"github.com/ohohleo/violin/opengl"
	"log"
	//"net/http"
)

func main() {

	accelerometer, err := input.AccelGyroSerial("/dev/ttyACM0", 38400)
	if err != nil {
		log.Fatal(err)
	}

	// err = api.New()

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = api.NewStream(accelerometer)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	//go func() {
	window, err := opengl.CreateWindow()
	if err != nil {
		panic(err)
	}

	object := window.AddObject()
	go func() {
		for {
			values := <-accelerometer
			fmt.Printf("%s", values)
			object.GetTransform().SetRotate(
				values.QuaternionW,
				values.QuaternionX,
				values.QuaternionY,
				values.QuaternionZ)
		}
	}()

	window.Start()

	//}()

	// log.Println("Listening :5000 ...")
	// http.ListenAndServe(":5000", nil)
}
