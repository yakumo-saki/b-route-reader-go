package bp35a1

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/yakumo-saki/b-route-reader-go/src/config"
	"github.com/yakumo-saki/b-route-reader-go/src/echonet"
	"github.com/yakumo-saki/b-route-reader-go/src/echonet/smartmeter"
)

var LAST_TID uint16 = 1000

var parser smartmeter.ELSmartMeterParser

// 一度だけ取得すればあとは変わらないものを取得して、パーサーに渡す
func GetSmartMeterInitialData(ipv6 string) error {
	var err error
	tid := echonet.TransactionId(9000)
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

	echonetResponse := findEchonetResponses(ret, tidStr)
	if len(echonetResponse) == 0 {
		return fmt.Errorf("failed to get echonet response: %w", err)
	}
	el, err := echonet.Parse(echonetResponse)
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

// あたらしいTransactionIDを生成して文字列で返す.
func getNewTid() echonet.TransactionId {
	LAST_TID = LAST_TID + 1
	if LAST_TID >= 9000 {
		LAST_TID = 1000
		log.Info().Msg("TID reached 9000. set back to 1000.")
	}
	return echonet.TransactionId(LAST_TID)

}

// 瞬間電力消費量を取得する
func GetNowConsumptionData(ipv6 string) (ElectricData, error) {

	tid := getNewTid()

	// 積算追加
	targets := []byte{
		echonet.P_NOW_DENRYOKU,
		echonet.P_NOW_DENRYUU,
	}
	msg := echonet.CreateEchonetGetMessage(tid, targets)

	electricData, err := querySmartMeter(ipv6, tid, msg)
	if err != nil {
		return electricData, err
	}

	return electricData, nil
}

// 瞬間電力消費量を取得する
func GetDeltaConsumptionData(ipv6 string) (ElectricData, error) {

	tid := getNewTid()

	// 積算追加
	targets := []byte{
		echonet.P_DELTA_DENRYOKU,
		echonet.P_DELTA_DENRYOKU_R,
	}
	msg := echonet.CreateEchonetGetMessage(tid, targets)

	electricData, err := querySmartMeter(ipv6, tid, msg)
	if err != nil {
		return electricData, fmt.Errorf("failed to get delta consumption: %w", err)
	}

	return electricData, nil
}

func querySmartMeter(ipv6 string, tid echonet.TransactionId, msg []byte) (ElectricData, error) {
	tidStr := fmt.Sprintf("%04d", tid)

	nullResult := ElectricData{}
	var err error

	ret, err := sendRequestWithRetry(ipv6, tid, msg)
	if err != nil {
		return nullResult, fmt.Errorf("error occured while sending request: %w", err)
	}

	elret := findEchonetResponses(ret, tidStr)
	if len(elret) == 0 {
		return nullResult, fmt.Errorf("echonet responses not found")
	}

	elmsg, err := echonet.Parse(elret)
	if err != nil {
		return nullResult, err
	}

	electricData, err := parseELMsgToElectricData(elmsg)
	if err != nil {
		return nullResult, err
	}

	return electricData, nil
}

func sendRequestWithRetry(ipv6 string, tid echonet.TransactionId, msg []byte) ([]string, error) {
	nullResult := []string{}
	tidStr := fmt.Sprintf("%04d", tid)

	retry := 0
	for { // retry loop
		err := skSendTo(ipv6, msg)
		if err != nil {
			return nullResult, err
		}

		log.Debug().Msgf("Echonet property GET message sent.")

		ret, err := waitForResultERXUDP(tidStr)
		if err != nil {
			if strings.Contains(err.Error(), "timeout reached") {
				if retry < config.MAX_ECHONET_GET_RETRY {
					log.Warn().Msgf("No smartmeter response. retrying %d/%d",
						(retry + 1), config.MAX_ECHONET_GET_RETRY)
					retry = retry + 1
					continue
				} else {
					// ここに来る原因は、長時間動かしている場合にPAN認証が期限切れになる場合がほとんど。
					// 本当はPAN認証からやり直せばよいのだが、実装するのも面倒なので異常終了としてsystemdに再起動してもらう
					log.Error().Msgf("Retry limit exceed. give up. needs restart")
					e := fmt.Errorf("no response. retry limit exceed")
					return nullResult, e
				}
			}

			return nullResult, err
		}

		if retry > 0 {
			log.Info().Msgf("Retry success. resuming normal operation.")
		}

		return ret, nil
	}
}

func parseELMsgToElectricData(elmsg echonet.EchonetLite) (ElectricData, error) {

	ret := ElectricData{}
	for key, value := range elmsg.Properties {
		switch key {
		case echonet.P_NOW_DENRYOKU:
			nowDenryoku, err := parser.ParseE7NowDenryoku(value)
			if err != nil {
				return ret, err
			}
			ret[fmt.Sprintf("%02X", key)] = float64(nowDenryoku)
		case echonet.P_DELTA_DENRYOKU:
			deltaDenryoku, err := parser.ParseE0DeltaDenryoku(value)
			if err != nil {
				return ret, err
			}
			ret[fmt.Sprintf("%02X", key)] = float64(deltaDenryoku)
		case echonet.P_DELTA_DENRYOKU_R:
			deltaDenryoku, err := parser.ParseE0DeltaDenryoku(value)
			if err != nil {
				return ret, err
			}
			ret[fmt.Sprintf("%02X", key)] = float64(deltaDenryoku)
		case echonet.P_NOW_DENRYUU:
			nowDenryuu, err := parser.ParseE8NowDenryuu(value)
			if err != nil {
				return ret, err
			}
			ret[fmt.Sprintf("%02X_Rphase", key)] = float64(nowDenryuu.Rphase)
			ret[fmt.Sprintf("%02X_Tphase", key)] = float64(nowDenryuu.Tphase)
			ret[fmt.Sprintf("%02X", key)] = float64(nowDenryuu.Total)
		}
	}

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
func findEchonetResponses(received []string, tid string) string {
	for _, v := range received {
		if isEchonetUnicastResponse(v, tid) {
			values := strings.Split(v, " ")
			return values[8]
		}
	}
	return ""
}

// ERXUDP FE80:0000:0000:0000:021C:6400:03CD:76A4 FE80:0000:0000:0000:021D:1290:0004:263D
// 0E1A 0E1A 001C640003CD76A4 1 0018 1081100202880105FF017202E7040000029CE80400280028
func isEchonetUnicastResponse(receivedErxUDP string, tid string) bool {
	if !strings.HasPrefix(receivedErxUDP, "ERXUDP") {
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
	case strings.HasPrefix(values[3], "02CC"):
		log.Debug().Msg("source port is 02CC(716:PANA) not echonet lite")
		return false
	case !strings.HasPrefix(values[3], "0E1A"):
		log.Debug().Msg("source port is not 0E1A(3610)")
		return false // 送信元ポートが3610
	case strings.HasPrefix(values[4], "02CC"):
		log.Debug().Msg("source port is 02CC(716:PANA) not echonet lite")
		return false
	case !strings.HasPrefix(values[4], "0E1A"):
		log.Debug().Msg("destination port is not 0E1A(3610)")
		return false // 送信先ポートが3610
	case !strings.Contains(values[8], tid):
		log.Debug().Msgf("TransactionId %s is not contains in echonet lite msg", tid)
		return false // TransactionIDがデータに含まれる
	}

	return true
}
