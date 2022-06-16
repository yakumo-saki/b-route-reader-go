package bp35a1

import (
	"github.com/yakumo-saki/b-route-reader-go/src/config"
)

func InitializeBrouteConnection() error {
	isAscii, err := isAsciiMode()
	if err != nil {
		return err
	}
	if isAscii == false {
		// WOPT 1
	}

	// ID PWD
	err = setBrouteId(config.B_ROUTE_ID)
	if err != nil {
		return err
	}
	err = setBroutePassword(config.B_ROUTE_PASSWORD)
	if err != nil {
		return err
	}

}
