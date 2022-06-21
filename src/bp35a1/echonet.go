package bp35a1

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/yakumo-saki/b-route-reader-go/src/echonet"
)

var UNIT byte
var LAST_TID uint16 = 1000

func InitEchonet(ipv6 string) error {
	var err error

	msg := echonet.CreateEchonetGetMessage(9000, []byte{echonet.P_UNIT})

	err = skSendTo(ipv6, msg)
	if err != nil {
		return err
	}

	log.Debug().Msgf("Echonet GET message sent.")

	_, err = waitForResultERXUDP()
	if err != nil {
		return err
	}

	return nil
}

//
func GetBrouteData(ipv6 string) ([]string, error) {

	var err error
	result := []string{"ERR"}

	tid := LAST_TID + 1
	if tid >= 9000 {
		LAST_TID = 1000
		tid = 1000
		log.Info().Msg("TID reached 9000. set back to 1000.")
	}

	msg := echonet.CreateEchonetGetMessage(tid, []byte{echonet.P_NOW_DENRYOKU, echonet.P_NOW_DENRYUU})

	err = skSendTo(ipv6, msg)
	if err != nil {
		return result, err
	}

	log.Debug().Msgf("Echonet property GET message sent.")

	ret, err := waitForResultERXUDP()
	if err != nil {
		return result, err
	}

	LAST_TID = tid

	// DEBUG: SKSENDTO のDATALENがバグっていればコマンドが欠けたりしてエラーになるはず
	// err = connectionTest()
	// if err != nil {
	// 	return result, err
	// }

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
