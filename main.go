package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
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
	config.Initialize()
	logger.Initiallize()

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

	for {

		ret, err := bp35a1.GetElectricData(ipv6)
		if err != nil {
			log.Err(err).Msg("Error occured while getting smartmeter data")
			exitcode = 1
			goto EXIT
		}

		log.Info().Msgf("Smartmeter Response: %s", ret)
		err = handleResult(ret)
		if err != nil {
			log.Err(err).Msg("Error occured while executing EXEC_CMD")
			exitcode = 1
			goto EXIT
		}

		// 連続でデータを取得しないためのwait。本当は規格的に20秒以上の間隔が必要
		log.Info().Msg("Wait for request data...")
		time.Sleep(10 * time.Second)
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

func handleResult(data bp35a1.ElectricData) error {

	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	f, err := ioutil.TempFile(os.TempDir(), "b-route-")
	if err != nil {
		return err
	}

	written, err := f.Write(json)
	if err != nil {
		return err
	}
	if written != len(json) {
		return fmt.Errorf("bytes written != actual")
	}

	filepath := f.Name()
	f.Close()

	// exec
	output, err := exec.Command(config.EXEC_CMD, f.Name()).CombinedOutput()
	if err != nil {
		return err
	}
	outputByteStringsToLog(output)

	os.Remove(filepath)

	return nil
}

func outputByteStringsToLog(byteStrings []byte) {
	newline := "\n"
	switch runtime.GOOS {
	case "windows":
		newline = "\r\n"
	case "darwin":
		newline = "\n"
	case "linux":
		newline = "\n"
	}

	allStrings := string(byteStrings)
	lines := strings.Split(allStrings, newline)
	for _, v := range lines {
		line := v
		line = strings.ReplaceAll(line, "\r", "")
		line = strings.ReplaceAll(line, "\n", "")
		if len(line) > 0 {
			log.Info().Msgf("EXEC OUTPUT: %s", line)
		}
	}
}
