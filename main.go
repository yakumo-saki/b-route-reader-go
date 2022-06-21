package main

import (
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/yakumo-saki/b-route-reader-go/src/bp35a1"
	"github.com/yakumo-saki/b-route-reader-go/src/config"
	"github.com/yakumo-saki/b-route-reader-go/src/logger"
)

var exitcode = 0

func main() {
	logger.Initiallize()
	config.Initialize()

	var ipv6 string

	log.Info().Msg("Start")
	err := bp35a1.Connect()
	if err != nil {
		log.Err(err).Msg("Serial port open error. Exiting.")
		exitcode = 1
		goto EXIT
	}

	err = bp35a1.StartConnection()
	if err != nil {
		log.Err(err).Msg("Test connection failed. Exiting.")
		exitcode = 1
		goto EXIT
	}

	ipv6, err = bp35a1.InitializeBrouteConnection()
	if err != nil {
		log.Err(err).Msg("Cannot initialize B-route connection. Exiting.")
		exitcode = 1
		goto EXIT
	}

	// echonet start
	bp35a1.InitEchonet(ipv6)

	// main loop
	for i := 0; i < 1; i++ {
		ret, err := bp35a1.GetBrouteData(ipv6)
		if err != nil {
			log.Err(err).Msg("Error occured while getting smartmeter data")
			exitcode = 1
			goto EXIT
		}

		log.Debug().Msgf("%s", ret)
		time.Sleep(10 * time.Second)
	}

EXIT:
	err = bp35a1.Close()
	if err != nil {
		log.Err(err).Msg("Error occured in close connection. do nothing.")
	}

	os.Exit(exitcode)
}
