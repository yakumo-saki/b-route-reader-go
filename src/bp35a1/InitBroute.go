package bp35a1

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/yakumo-saki/b-route-reader-go/src/config"
)

func InitializeBrouteConnection() error {
	isAscii, err := isAsciiMode()
	if err != nil {
		return err
	}
	if !isAscii {
		// WOPT 1
		log.Warn().Msg("WOPT 1 is not implemented. maybe not working.")
		return fmt.Errorf("command 'WOPT 1' is not implemented")
	}

	err = setupIdAndPassword()
	if err != nil {
		return err
	}

	sm, err := searchSmartMeter()
	if err != nil {
		return err
	}

	log.Info().Msgf("Found smartmeter %s", sm)

	ipv6, err := convertPanIdToIpv6(sm.Addr)
	if err != nil {
		return err
	}

	log.Info().Msgf("Smartmeter address is %s", ipv6)

	return nil
}

func setupIdAndPassword() error {
	// ID PWD
	err := setBrouteId(config.B_ROUTE_ID)
	if err != nil {
		return err
	}

	err = setBroutePassword(config.B_ROUTE_PASSWORD)
	if err != nil {
		return err
	}

	return nil
}

func searchSmartMeter() (SmartMeter, error) {

	log.Info().Msg("Active scan start")

	sm, err := activeScan()
	if err != nil {
		return sm, err
	}

	return sm, err
}
