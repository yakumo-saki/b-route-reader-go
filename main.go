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

	log.Info().Msg("Start")
	err := bp35a1.Connect()
	if err != nil {
		log.Err(err).Msg("Serial port open error. Exiting.")
		exitcode = 1
		goto EXIT
	}

	err = runWithSerialPort()
	if err != nil {
		exitcode = 1
		log.Err(err).Msg("ERR")
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

func runWithSerialPort() error {
	var err error
	var ipv6 string

	err = bp35a1.StartConnection()
	if err != nil {
		return fmt.Errorf("test connection failed: %w", err)
	}

	ipv6, err = bp35a1.InitializeBrouteConnection()
	if err != nil {
		log.Err(err).Msg(". Exiting.")
		return fmt.Errorf("cannot initialize B-route connection: %w", err)
	}

	// echonet start
	err = bp35a1.GetSmartMeterInitialData(ipv6)
	if err != nil {
		return fmt.Errorf("error occured while initializing echonet lite: %w", err)
	}

	log.Info().Msg("Starting main loop")

	// TODO シグナルハンドリング
	// TODO 積算電力量を数分に一回取得する 3min ?

	for {

		ret, err := bp35a1.GetElectricData(ipv6)
		if err != nil {
			return fmt.Errorf("error occured while getting smartmeter data: %w", err)
		}

		log.Info().Msgf("Smartmeter Response: %v", ret)
		err = handleResult(ret)
		if err != nil {
			return fmt.Errorf("error occured while executing %s: %w", config.EXEC_CMD, err)
		}

		// 連続でデータを取得しないためのwait。本当は規格的に20秒以上の間隔が必要
		log.Info().Msg("Wait for request data...")
		time.Sleep(10 * time.Second)
	}

	return nil
}

func handleResult(data bp35a1.ElectricData) error {

	jsonMap := map[string] interface{}
	for k, v := range data {
		jsonMap[k] = v
	}
	jsonMap["datetime"] = time.Now().Format("2006-01-02T15:04:05.999Z")

	json, err := json.Marshal(jsonMap)
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
