package bp35a1

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

func isAsciiMode() (bool, error) {
	sendCommand("ROPT")

	responses, err := waitForResult()
	if err != nil {
		return false, err
	}

	for _, v := range responses {
		if strings.HasPrefix(v, "OK 01") {
			log.Debug().Msg("ROPT OK. Already ASCII mode.")
			return true, nil
		}
	}

	log.Debug().Msg("ROPT OK. BINARY mode.")
	return false, nil
}

func setBrouteId(id string) error {
	err := sendCommand(fmt.Sprintf("SKSETRBID %s", id))
	return err
}

func setBroutePassword(password string) error {
	err := sendCommand(fmt.Sprintf("SKSETPWD C %s", password))
	return err
}

func activeScan() error {
	SKSCAN 2 FFFFFFFF 6
}