package input

import (
	"bufio"
	"fmt"
	"github.com/tarm/serial"
	"log"
	"strconv"
	"strings"
)

type AccelGyro struct {
	aX int
	aY int
	aZ int

	gX int
	gY int
	gZ int
}

func (a *AccelGyro) String() string {
	return fmt.Sprintf("x%dy%dz%d", a.aX, a.aY, a.aZ)
}

// Etablie la connexion avec le port série spécifié pour récupérer les
// données provenant de l'accéléromètre & du gyroscope
func AccelGyroSerial(device string, baudrate int, isReadable bool) (chan AccelGyro, error) {

	// Création du channel
	channel := make(chan AccelGyro)

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

	if isReadable {
		go fromReadable(buf, channel)
	} else {
		go fromBynary(buf, channel)
	}

	return channel, nil
}

// Réception des valeurs lisibles OUTPUT_READABLE_ACCELGYRO
func fromReadable(buf *bufio.Reader, channel chan AccelGyro) {

	values := make([]int, 6)

	for {
		rcv, err := buf.ReadBytes('\n')
		if err != nil {
			return
		}

		strValues := strings.SplitN(
			strings.TrimSpace(string(rcv)), "\t", 6)

		if len(strValues) != 6 {
			continue
		}

		for idx, strVal := range strValues {
			number, err := strconv.Atoi(strVal)
			if err != nil {
				log.Printf("Unexpected value at id #%d (%s): %s",
					idx, strVal, err.Error())
				continue
			}

			values[idx] = number
		}

		accelerometer := AccelGyro{
			aX: values[0],
			aY: values[1],
			aZ: values[2],
			gX: values[3],
			gY: values[4],
			gZ: values[5],
		}

		log.Printf("%+v", accelerometer)
	}
}

// Réception des valeurs binaires OUTPUT_BINARY_ACCELGYRO
func fromBynary(buf *bufio.Reader, channel chan AccelGyro) {
	var previousBuffer []byte
	var previousSize int
	var crc int

	for {
		rcv, err := buf.ReadBytes('\n')
		if err != nil {
			return
		}

		// Taille du buffer trop petite
		if len(rcv) < 17 {

			// Stockage du buffer le temps de recevoir le reste
			if previousBuffer == nil {

				previousBuffer = make([]byte, 17)
				previousSize = len(rcv)

				for i := 0; i < len(rcv); i++ {
					previousBuffer[i] = rcv[i]
				}

				continue
			}

			if previousSize+len(rcv) != 17 {
				goto error
			}

			// sinon on rajoute ce que l'on vient d'obtenir au
			// buffer
			for i := 0; i < len(rcv); i++ {
				previousBuffer[previousSize+i] = rcv[i]
			}

			rcv = previousBuffer
		}

		previousBuffer = nil
		previousSize = 0

		// Vérification de la taille et de l'entête
		if len(rcv) == 17 && rcv[0] == ':' {

			// Récupération de la taille
			len := rcv[1]

			// Récupération du crc
			crc = uint16ToInt(rcv[len+2], rcv[len+3])

			rcv = rcv[2 : len+2]

			// Vérification du crc
			expect := crc16(rcv)
			if crc != expect {
				log.Printf("invalid crc: got %04X, expect %04X\n", crc, expect)
				goto error
			}

			accelerometer := AccelGyro{
				aX: uint16ToInt16(rcv[0], rcv[1]),
				aY: uint16ToInt16(rcv[2], rcv[3]),
				aZ: uint16ToInt16(rcv[4], rcv[5]),
				gX: uint16ToInt16(rcv[6], rcv[7]),
				gY: uint16ToInt16(rcv[8], rcv[9]),
				gZ: uint16ToInt16(rcv[10], rcv[11]),
			}

			log.Printf("%+v", accelerometer)

			continue

		}

	error:
		previousBuffer = nil
		previousSize = 0
		log.Printf("invalid data received size:%d value:%+v\n", len(rcv), rcv)
		continue
	}
}

func uint16ToInt(hi byte, low byte) int {
	return int(hi)<<8 + int(low)
}

func uint16ToInt16(hi byte, low byte) int {
	value := uint16ToInt(hi, low)
	if value&0x4000 > 0 {
		return -1 * (-value & 0x7fff)
	}

	return value
}

// Implements CRC-CCITT (Kermit)
func crc16(buf []byte) int {

	var crc uint16
	var polynom uint16 = 0x1081

	for i := 0; i < len(buf); i++ {
		b := uint16(buf[i])
		q := (crc ^ b) & 0x0f
		crc = (crc >> 4) ^ (q * polynom)
		q = (crc ^ (b >> 4)) & 0xf
		crc = (crc >> 4) ^ (q * polynom)
	}

	return int((crc >> 8) ^ (crc << 8))
}
