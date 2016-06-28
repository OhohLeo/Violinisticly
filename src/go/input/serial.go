package input

import (
	"bufio"
	"fmt"
	"github.com/tarm/serial"
	"log"
	"regexp"
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
	valid := regexp.MustCompile(`^[0-9]+-[0-9]+-[0-9]+`)

	go func() {

		for {

			values, err := buf.ReadBytes('\n')

			if err != nil {
				return
			}

			if valid.Match(values) {

				// Retire les '\r\n' finaux
				values = values[:len(values)-2]

				// Récupération des valeurs X, Y & Z
				split := strings.SplitN(string(values), "-", 3)

				x, err := strconv.Atoi(split[0])
				if err != nil {
					log.Printf("x %+v", err)
					return
				}

				y, err := strconv.Atoi(split[1])
				if err != nil {
					log.Printf("y %+v", err)
					return
				}

				z, err := strconv.Atoi(split[2])
				if err != nil {
					log.Printf("z %+v", err)
					return
				}

				accelerometer := Accelerometer{
					X: x,
					Y: y,
					Z: z,
				}

				channel <- accelerometer
			}
		}
	}()

	return channel, nil
}
