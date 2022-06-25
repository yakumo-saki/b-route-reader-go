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
	ret := run()
	os.Exit(ret)
}

func run() int {
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
	err = bp35a1.GetSmartMeterInitialData(ipv6)
	if err != nil {
		log.Err(err).Msg("Error occured while initializing echonet lite")
		exitcode = 1
		goto EXIT
	}

	// main loop
	for i := 0; i < 2; i++ {
		ret, err := bp35a1.GetElectricData(ipv6)
		if err != nil {
			log.Err(err).Msg("Error occured while getting smartmeter data")
			exitcode = 1
			goto EXIT
		}

		// TODO implement this!
		log.Info().Msgf("%s", ret)

		// 連続でデータを取得しないためのwait。本当は規格的に20秒以上の間隔が必要
		if i > 1 {
			log.Info().Msg("Wait for request data...")
			time.Sleep(10 * time.Second)
		}
	}

EXIT:
	err = bp35a1.Close()
	if err != nil {
		log.Err(err).Msg("Error occured in close connection. do nothing.")
	}

	if exitcode == 0 {
		log.Info().Msg("Normal end.")
	}

	return exitcode
}
