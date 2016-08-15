package input

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/tarm/serial"
	"log"
	"math"
)

const (
	QUATERNION = 1 << iota
	EULER
	YAWPITCHROLL
	REALACCEL
	WORLDACCEL

	IDX_HEADER = 0
	IDX_LEN    = 1

	SIZE_HEADER = 2
	SIZE_CRC    = 2
)

var INIT_STATUS []string = []string{
	"MPU init",
	"MPU connection",
	"DMP init",
	"DMP interrupt status",
	"FIFO overflow!",
}

type AccelGyro struct {
	Status int

	// Quaternion
	QuaternionW float32
	QuaternionX float32
	QuaternionY float32
	QuaternionZ float32

	// Euler
	EulerX float32
	EulerY float32
	EulerZ float32

	// Yaw/Pitch/Roll
	Yaw   float32
	Pitch float32
	Roll  float32

	// Real Acceleration
	RealX float32
	RealY float32
	RealZ float32

	// World Acceleration
	WorldX float32
	WorldY float32
	WorldZ float32
}

func (a *AccelGyro) String() string {

	result := ""

	if (a.Status & QUATERNION) > 0 {
		result += fmt.Sprintf("quaternion:\tw:%f\tx:%f\ty:%f\tz:%f",
			a.QuaternionW, a.QuaternionX, a.QuaternionY, a.QuaternionZ)
	}

	if (a.Status & EULER) > 0 {
		result += fmt.Sprintf("euler:\tx:%f\ty:%f\tz:%f",
			a.EulerX, a.EulerY, a.EulerZ)
	}

	if (a.Status & YAWPITCHROLL) > 0 {
		result += fmt.Sprintf("yaw/pitch/roll:\tyaw:%f\tpitch:%f\troll:%f",
			a.Yaw, a.Pitch, a.Roll)
	}

	if (a.Status & REALACCEL) > 0 {
		result += fmt.Sprintf("real:\tx:%f\ty:%f\tz:%f",
			a.RealX, a.RealY, a.RealZ)
	}

	if (a.Status & WORLDACCEL) > 0 {
		result += fmt.Sprintf("world:\tx:%f\ty:%f\tz:%f",
			a.WorldX, a.WorldY, a.WorldZ)
	}

	return result
}

// Etablie la connexion avec le port série spécifié pour récupérer les
// données provenant de l'accéléromètre & du gyroscope
func AccelGyroSerial(device string, baudrate int) (chan *AccelGyro, error) {

	// Création du channel
	channel := make(chan *AccelGyro)

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

	go fromBinary(buf, channel)

	return channel, nil
}

var previousBuffer []byte

// Réception des valeurs binaires OUTPUT_BINARY_ACCELGYRO
func fromBinary(buf *bufio.Reader, channel chan *AccelGyro) {
	var crc int

	for {
		rcv, err := buf.ReadBytes('\n')
		if err != nil {
			return
		}

		// Un buffer est déjà stocké : on concatène les nouvelles données reçues
		if previousBuffer != nil {
			previousBuffer = append(previousBuffer, rcv...)
			rcv = previousBuffer
		}

		// Taille du buffer trop petite
		if len(rcv) < 7 {
			continue
		}

		// Vérification du caractère d'entête ':'
		if rcv[IDX_HEADER] == ':' {

			// Récupération de la taille
			length := int(rcv[IDX_LEN])

			// Vérification de la taille + CRC (2 octets)
			if len(rcv) < IDX_LEN+length+SIZE_CRC+1 {
				previousBuffer = rcv
				continue
			}

			// Trame dont la taille est correcte : réinitialisation du buffer
			previousBuffer = nil

			// Récupération du crc
			crc = uint16ToInt(rcv[SIZE_HEADER+length], rcv[SIZE_HEADER+length+1])

			// Suppression de l'entête et du CRC
			rcv = rcv[SIZE_HEADER : SIZE_HEADER+length]

			// Vérification du crc
			expect := crc16(rcv)
			if crc != expect {
				log.Printf("invalid crc: got %04X, expect %04X\n", crc, expect)
				continue
			}

			// Récupération du status
			status := rcv[0]

			// Récupération d'un status d'initialisation
			if length == 2 && status < 5 {
				log.Printf("%s: %d\n", INIT_STATUS[status], rcv[1])
				continue
			}

			rcv = rcv[1:]

			values := &AccelGyro{
				Status: int(status),
			}

			if status&QUATERNION > 0 {

				if len(rcv) < 16 {
					log.Printf("Quaternion: invalid got %d expect %d\n", len(rcv), 16)
					continue
				}

				values.QuaternionW = float32frombytes(rcv)
				values.QuaternionX = float32frombytes(rcv[4:])
				values.QuaternionY = float32frombytes(rcv[8:])
				values.QuaternionZ = float32frombytes(rcv[12:])

				rcv = rcv[16:]
			}

			if status&EULER > 0 {

				if len(rcv) < 12 {
					log.Printf("Euler: invalid got %d expect %d\n", len(rcv), 12)
					continue
				}

				values.EulerX = float32frombytes(rcv)
				values.EulerY = float32frombytes(rcv[4:])
				values.EulerZ = float32frombytes(rcv[8:])

				rcv = rcv[12:]
			}

			if status&YAWPITCHROLL > 0 {

				if len(rcv) < 12 {
					log.Printf("Euler: invalid got %d expect %d\n", len(rcv), 12)
					continue
				}

				values.Yaw = float32frombytes(rcv)
				values.Pitch = float32frombytes(rcv[4:])
				values.Roll = float32frombytes(rcv[8:])

				rcv = rcv[12:]
			}

			if status&REALACCEL > 0 {

				if len(rcv) < 12 {
					log.Printf("Euler: invalid got %d expect %d\n", len(rcv), 12)
					continue
				}

				values.RealX = float32frombytes(rcv)
				values.RealY = float32frombytes(rcv[4:])
				values.RealZ = float32frombytes(rcv[8:])

				rcv = rcv[12:]
			}

			if status&WORLDACCEL > 0 {

				if len(rcv) < 12 {
					log.Printf("Euler: invalid got %d expect %d\n", len(rcv), 12)
					continue
				}

				values.WorldX = float32frombytes(rcv)
				values.WorldY = float32frombytes(rcv[4:])
				values.WorldZ = float32frombytes(rcv[8:])

				rcv = rcv[12:]
			}

			channel <- values
		}
	}
}

func float32frombytes(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

func uint16ToInt(hi byte, low byte) int {
	return int(hi)<<8 + int(low)
}

func uint16ToInt16(hi byte, low byte) int64 {
	value := uint16ToInt(hi, low)
	if value&0x4000 > 0 {
		return int64(-1 * (-value & 0x7fff))
	}

	return int64(value)
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
