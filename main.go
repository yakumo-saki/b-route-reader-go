package main

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/yakumo-saki/b-route-reader-go/src/bp35a1"
	"github.com/yakumo-saki/b-route-reader-go/src/config"
	"github.com/yakumo-saki/b-route-reader-go/src/logger"
)

var exitcode = 0

func main() {
	logger.Initiallize()
	config.Initialize()

	log.Info().Msg("Start")
	err := bp35a1.Connect()
	if err != nil {
		log.Err(err).Msg("Serial port open error. Exiting.")
		exitcode = 1
		goto EXIT
	}

	err = bp35a1.TestConnection()
	if err != nil {
		log.Err(err).Msg("Serial port opened. Exiting.")
		exitcode = 1
		goto EXIT
	}

	err = bp35a1.InitializeBrouteConnection()
	if err != nil {
		log.Err(err).Msg("Serial port opened. Exiting.")
		exitcode = 1
		goto EXIT
	}

	// main loop

EXIT:
	err = bp35a1.Close()
	if err != nil {
		log.Err(err).Msg("Error occured in close connection. do nothing.")
	}

	os.Exit(exitcode)
}
