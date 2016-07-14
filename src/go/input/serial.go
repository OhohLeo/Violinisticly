package input

import (
	"bufio"
	"fmt"
	"github.com/tarm/serial"
	"log"
	"strconv"
	"strings"
)

type Accelerometer struct {
	X int
	Y int
	Z int
}

func (a *Accelerometer) String() string {
	return fmt.Sprintf("x%dy%dz%d", a.X, a.Y, a.Z)
}

func AccelerometerSerial(device string, baudrate int) (chan Accelerometer, error) {

	// Création du channel
	channel := make(chan Accelerometer)

	// Configuration du port série
	c := &serial.Config{
		Name: device,
		Baud: baudrate,
	}

	// Ouverture
	s, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}

	// Lecture
	buf := bufio.NewReader(s)

	go func() {

		for {
			rcv, err := buf.ReadBytes('\n')
			if err != nil {
				return
			}

			s := strings.TrimSpace(string(rcv))

			values := strings.SplitN(s, "-", 3)

			x, err := strconv.Atoi(values[0])
			if err != nil {
				log.Printf("x %+v", err)
				return
			}

			y, err := strconv.Atoi(values[1])
			if err != nil {
				log.Printf("y %+v", err)
				return
			}

			z, err := strconv.Atoi(values[2])
			if err != nil {
				log.Printf("z %+v", err)
				return
			}

			accelerometer := Accelerometer{
				X: x,
				Y: y,
				Z: z,
			}

			log.Printf("%+v", accelerometer)

			channel <- accelerometer
		}
	}()

	return channel, nil
}
