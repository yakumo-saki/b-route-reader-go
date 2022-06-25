package bp35a1

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/yakumo-saki/b-route-reader-go/src/echonet"
	"github.com/yakumo-saki/b-route-reader-go/src/echonet/smartmeter"
)

var LAST_TID uint16 = 1000

var parser smartmeter.ELSmartMeterParser

// 一度だけ取得すればあとは変わらないものを取得して、パーサーに渡す
func GetSmartMeterInitialData(ipv6 string) error {
	var err error
	tid := uint16(9000)
	msg := echonet.CreateEchonetGetMessage(tid, []byte{echonet.P_DELTA_UNIT, echonet.P_MULTIPLIER})

	err = skSendTo(ipv6, msg)
	if err != nil {
		return err
	}

	log.Debug().Msgf("Echonet GET message sent.")
	tidStr := fmt.Sprintf("%04d", tid)

	ret, err := waitForResultERXUDP(tidStr)
	if err != nil {
		return err
	}

	responses := findEchonetResponse(ret, tidStr)
	if len(responses) == 0 {
		return fmt.Errorf("failed to get echonet response: %w", err)
	}
	el, err := echonet.Parse(responses[0])
	if err != nil {
		return fmt.Errorf("failed to parse response as echonet msg: %w", err)
	}

	parser = smartmeter.ELSmartMeterParser{}
	_, err = parser.ParseAndStoreE1DeltaUnit(el.Properties[echonet.P_DELTA_UNIT])
	if err != nil {
		return fmt.Errorf("failed to parse P_DELTA_UNIT(D0): %w", err)
	}

	_, err = parser.ParseAndStoreD3Multiplier(el.Properties[echonet.P_MULTIPLIER])
	if err != nil {
		return fmt.Errorf("failed to parse P_MULTIPLIER(D3): %w", err)
	}

	return nil
}

//
func GetElectricData(ipv6 string) (ElectricData, error) {

	var err error
	nullResult := ElectricData{}

	tid := LAST_TID + 1
	if tid >= 9000 {
		LAST_TID = 1000
		tid = 1000
		log.Info().Msg("TID reached 9000. set back to 1000.")
	}

	// 積算追加
	targets := []byte{
		echonet.P_NOW_DENRYOKU,
		echonet.P_NOW_DENRYUU,
		echonet.P_DELTA_DENRYOKU,
	}
	msg := echonet.CreateEchonetGetMessage(tid, targets)

	err = skSendTo(ipv6, msg)
	if err != nil {
		return nullResult, err
	}

	log.Debug().Msgf("Echonet property GET message sent.")

	tidStr := fmt.Sprintf("%04d", tid)
	ret, err := waitForResultERXUDP(tidStr)
	if err != nil {
		return nullResult, err
	}

	elret := findEchonetResponse(ret, tidStr)
	if len(elret) != 1 {
		return nullResult, fmt.Errorf("multiple echonet responses found, maybe bug")
	}

	elstr := elret[0]
	elmsg, err := echonet.Parse(elstr)
	if err != nil {
		return nullResult, err
	}

	electricData, err := parseELMsgToElectricData(elmsg)
	if err != nil {
		return nullResult, err
	}

	LAST_TID = tid

	return electricData, nil
}

func parseELMsgToElectricData(elmsg echonet.EchonetLite) (ElectricData, error) {

	ret := ElectricData{}
	nowDenryokuBytes := elmsg.Properties[echonet.P_NOW_DENRYOKU]
	nowDenryoku, err := parser.ParseE7NowDenryoku(nowDenryokuBytes)
	if err != nil {
		return ret, err
	}

	deltaDenryokuBytes := elmsg.Properties[echonet.P_DELTA_DENRYOKU]
	deltaDenryoku, err := parser.ParseE0DeltaDenryoku(deltaDenryokuBytes)
	if err != nil {
		return ret, err
	}

	nowDenryuuBytes := elmsg.Properties[echonet.P_NOW_DENRYUU]
	nowDenryuu, err := parser.ParseE8NowDenryuu(nowDenryuuBytes)
	if err != nil {
		return ret, err
	}

	ret.DeltakWh = deltaDenryoku
	ret.Watt = nowDenryoku
	ret.RphaseAmp = nowDenryuu.Rphase
	ret.TphaseAmp = nowDenryuu.Tphase
	ret.TotalAmp = nowDenryuu.Total

	return ret, nil

}

// SKSENDTO コマンドを発行する。 Echonet通信用
// TEST
func skSendTo(ipv6 string, data []byte) error {

	secured := "1" // 1=YES 0=NO 受信したIPパケットを構成するMACフレームを暗号化するか（原文ママ）

	cmd := fmt.Sprintf("SKSENDTO 1 %s 0E1A %s %04X", ipv6, secured, len(data)) // DATALENは16進数
	cmd = cmd + " "                                                            // DATALENのあとにスペースを入れないとデータ入力待ちにならない
	n, err := port.Write([]byte(cmd))
	if err != nil {
		log.Err(err).Msgf("skSendTo: command send error")
		return err
	}
	if n != len(cmd) {
		return fmt.Errorf("skSendTo: unexpected command sent bytes. expected=%d actual=%d", len(cmd), n)
	}

	// データ入力待ちになるまでちょっと間隔を空ける
	time.Sleep(200 * time.Millisecond)

	n, err = port.Write(data)
	if err != nil {
		log.Err(err).Msgf("skSendTo: Data send error")
		return err
	}
	if n != len(data) {
		return fmt.Errorf("skSendTo: unexpected data sent bytes. expected=%d actual=%d", len(data), n)
	}

	log.Debug().Msgf("--> %sbinary<%s>", cmd, hex.EncodeToString(data))

	return nil
}

// Unicastで返ってくるEchonet Lite応答を抜き出して返す
func findEchonetResponse(received []string, tid string) []string {
	ret := make([]string, 0)
	for _, v := range received {
		if strings.Contains(v, tid) {
			ret = append(ret, v)
		}
	}
	return ret
}

// ERXUDP FE80:0000:0000:0000:021C:6400:03CD:76A4 FE80:0000:0000:0000:021D:1290:0004:263D
// 0E1A 0E1A 001C640003CD76A4 1 0018 1081100202880105FF017202E7040000029CE80400280028
func isEchonetUnicastResponse(receivedErxUDP string, tid string) bool {
	if !strings.HasPrefix(receivedErxUDP, "ERXUDP") {
		log.Debug().Msg("not start with ERXUDP")
		return false
	}
	values := strings.Split(receivedErxUDP, " ")
	if len(values) != 9 {
		log.Debug().Msg("ERXUDP splitted length != 9")
		return false
	}

	switch {
	case !strings.HasPrefix(values[1], "FE80"):
		log.Debug().Msg("source ipv6 address not begin with FE80")
		return false // 送信元がIPv6 local
	case !strings.HasPrefix(values[2], "FE80"):
		log.Debug().Msg("destination ipv6 address not begin with FE80")
		return false // 送信先がIPv6 local
	case !strings.HasPrefix(values[3], "0E1A"):
		log.Debug().Msg("source port is not 0E1A(3610)")
		return false // 送信元ポートが3610
	case !strings.HasPrefix(values[4], "0E1A"):
		log.Debug().Msg("destination port is not 0E1A(3610)")
		return false // 送信先ポートが3610
	case !strings.Contains(values[8], tid):
		log.Debug().Msgf("TransactionId %s is not contains in echonet lite msg", tid)
		return false // TransactionIDがデータに含まれる
	}

	return true
}
