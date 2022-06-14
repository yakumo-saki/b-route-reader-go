package bp35a1

import (
	"github.com/rs/zerolog/log"

	"go.bug.st/serial"
)

var port serial.Port

func Connect() error {
	port, err := serial.Open("/dev/ttyAMA0", &serial.Mode{})
	if err != nil {
		log.Err(err)
		return err
	}

	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	if err := port.SetMode(mode); err != nil {
		log.Err(err)
		return err
	}
	log.Info().Msg("Port open successs")

}
