package bp35a1

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

func sendReset() error {
	err := sendCommand("SKRESET")
	if err != nil {
		log.Error().Msg("Send SKRESET command fail.")
		return err
	}

	ret, err := waitForResult()
	if err != nil {
		return err
	}

	if !containsInResult(ret, RET_OK) {
		return fmt.Errorf("response is not OK")
	}

	return nil
}

// 接続が正しく確立して、コマンドが通るかテストする
// return がnilならOK
func connectionTest() error {
	err := sendCommand("SKVER")
	if err != nil {
		return err
	}

	_, err = waitForResult()
	if err != nil {
		return err
	}

	return nil
}

func isAsciiMode() (bool, error) {
	err := sendCommand("ROPT")
	if err != nil {
		return false, err
	}

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
	if err != nil {
		return err
	}

	ret, err := waitForResult()
	if endWithResult(ret, RET_OK) {
		return nil
	}

	return fmt.Errorf("response is not %s", RET_OK)
}

func setBroutePassword(password string) error {
	err := sendCommand(fmt.Sprintf("SKSETPWD C %s", password))
	if err != nil {
		return err
	}

	ret, err := waitForResult()
	if endWithResult(ret, RET_OK) {
		return nil
	}

	return fmt.Errorf("response is not %s", RET_OK)
}

func activeScan() (SmartMeter, error) {
	meter := SmartMeter{}

	for i := 0; i < 10; i++ {
		err := sendCommand("SKSCAN 2 FFFFFFFF 6")
		if err != nil {
			return meter, err
		}

		log.Debug().Msgf("wait for result...")

		ret, err := waitForResult()
		if err != nil {
			return meter, err
		}

		if containsInResult(ret, RET_SCAN_FOUND) {
			meter := parseActiveScanResult(ret)
			return meter, nil
		}
		if endWithResult(ret, RET_SCAN_COMPLETE) {
			break
		}

		log.Debug().Msgf("SmartMeter not found. retry (%d/%d).", i, 10)

	}

	return meter, fmt.Errorf("no SmartMeter found")
}

func parseActiveScanResult(ret []string) SmartMeter {
	sm := SmartMeter{}

	for _, v := range ret {
		switch {
		case strings.Contains(v, RET_PAN_CHANNEL):
			sm.Channel = parseActiveScanResultOne(v)
		case strings.Contains(v, RET_PAN_CHANNEL_PAGE):
			sm.ChannelPage = parseActiveScanResultOne(v)
		case strings.Contains(v, RET_PAN_ID):
			sm.PanId = parseActiveScanResultOne(v)
		case strings.Contains(v, RET_PAN_ADDR):
			sm.Addr = parseActiveScanResultOne(v)
		case strings.Contains(v, RET_PAN_LQI):
			sm.LQI = parseActiveScanResultOne(v)
		case strings.Contains(v, RET_PAN_PAIR_ID):
			sm.PairId = parseActiveScanResultOne(v)
		}
	}

	return sm
}

func parseActiveScanResultOne(colonSeparatedValue string) string {
	splitted := strings.SplitN(colonSeparatedValue, ":", 2)
	if len(splitted) != 2 {
		log.Error().Msgf("Unexpected EPANDESC return: %s", colonSeparatedValue)
	}

	return splitted[1]
}
