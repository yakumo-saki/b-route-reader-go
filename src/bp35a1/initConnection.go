package bp35a1

import (
	"github.com/rs/zerolog/log"
	"github.com/yakumo-saki/b-route-reader-go/src/config"

	"go.bug.st/serial"
)

var port serial.Port

const CRLF = "\r\n"

func Connect() error {
	var err error

	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	port, err = serial.Open(config.SERIAL, mode)
	if err != nil {
		log.Err(err).Str("PORT", config.SERIAL).Msg("Serial port open error")
		return err
	}

	if err := port.SetMode(mode); err != nil {
		log.Err(err).Msg("Serial port set mode error")
		return err
	}

	log.Debug().Msg("Port open successs")
	return nil
}

// Closeはgo-serial内でポートが開いているかの判定をしているので常に実行しても安全。
func Close() error {
	err := port.Close()
	if err == nil {
		log.Debug().Msg("Port close successs")
	}
	return err
}

func StartConnection() error {

	// すでに存在している応答があれば読み捨てる
	port.ResetInputBuffer()
	port.ResetOutputBuffer()

	err := sendReset()
	if err != nil {
		log.Err(err).Msg("Send reset command error")
		return err
	}

	err = setLocalEcho(false)
	if err != nil {
		log.Err(err).Msg("Stop echo command error")
		return err
	}

	err = connectionTest()
	if err != nil {
		log.Err(err).Msg("Connection test command error")
		return err
	}

	return nil
}
