package bp35a1

import (
	"bytes"
	"time"

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

func TestConnection() error {
	sendCommand("SKVER")

	waitForResult()

	log.Debug().Msg("SKVER OK")

	return nil
}

func waitForResult() ([]string, error) {

	log.Debug().Msg("Response start")
	BYTE_CR := []byte("\r")
	BYTE_LF := []byte("\n")

	port.SetReadTimeout(300 * time.Millisecond)

	var result []string

	// Read and print the response
	var byteBuf []byte
	readBuf := make([]byte, 100)
	for {
		// Reads up to 100 bytes
		n, err := port.Read(readBuf)
		if err != nil {
			log.Err(err).Msg("Read error")
			return nil, err
		}
		if n == 0 && len(byteBuf) == 0 {
			if len(result) > 0 {
				break
			}
		}

		byteBuf = append(byteBuf, readBuf[:n]...)

		if bytes.Contains(byteBuf, BYTE_CR) {
			crPos := bytes.Index(byteBuf, BYTE_CR)
			lineBuf := byteBuf[:crPos]

			byteBuf = byteBuf[crPos+len(BYTE_CR):] // 改行コードは削除

			// CRLFで区切られている場合、LFが残るので削除
			if bytes.HasPrefix(byteBuf, BYTE_LF) {
				byteBuf = byteBuf[1:]
				log.Debug().Msgf("<-- %s<CRLF>", string(lineBuf))
			} else {
				log.Debug().Msgf("<-- %s<CR>", string(lineBuf))
			}

			result = append(result, string(lineBuf))

			if crPos == 0 {
				break
			}
		}
	}

	log.Debug().Msg("Response done")
	return result, nil
}
