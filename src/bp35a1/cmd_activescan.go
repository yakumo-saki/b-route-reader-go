package bp35a1

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/yakumo-saki/b-route-reader-go/src/config"
)

func activeScan() (SmartMeter, error) {
	meter := SmartMeter{}

	for i := 0; i < config.ACTIVE_SCAN_COUNT; i++ {
		err := sendCommand("SKSCAN 2 FFFFFFFF 6")
		if err != nil {
			return meter, err
		}

		log.Debug().Msgf("wait for result...")

		ret, err := waitForResultSKSCAN()
		if err != nil {
			return meter, err
		}

		if containsInResult(ret, RET_SCAN_FOUND) {
			meter := parseActiveScanResult(ret)
			return meter, nil
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
