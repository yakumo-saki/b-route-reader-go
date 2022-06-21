package bp35a1

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/yakumo-saki/b-route-reader-go/src/echonet"
)

func InitEchonet(ipv6 string) error {
	var err error

	msg := echonet.CreateEchonetInfMessage(1234, []string{})
	log.Debug().Msgf("Echonet msg created -> %s", hex.EncodeToString(msg))

	err = skSendTo(ipv6, msg)
	if err != nil {
		return err
	}

	ret, err := waitForResultImpl([]string{"ERXUDP"})
	if err != nil {
		return err
	}
	dumpResult(ret)

	return nil
}

//
func GetBrouteData(ipv6 string) ([]string, error) {

	var err error
	result := []string{"ERR"}

	msg := echonet.CreateEchonetGetMessage(1234, []string{})
	log.Debug().Msgf("Echonet msg created -> %s", hex.EncodeToString(msg))

	err = skSendTo(ipv6, msg)
	if err != nil {
		return result, err
	}

	ret, err := waitForResultImpl([]string{"ERXUDP"})
	if err != nil {
		return result, err
	}
	dumpResult(ret)

	// ダミーでSKVERを実行して正しくSKSENDTOがされているかチェック
	// バグっていればSKVERというコマンドが欠けたりしてエラーになるはず
	err = sendCommand("SKVER")
	if err != nil {
		return result, err
	}

	_, err = waitForResultImpl([]string{})
	if err != nil {
		return result, err
	}

	return ret, nil
}

// SKSENDTO コマンドを発行する。 Echonet通信用
// TEST
func skSendTo(ipv6 string, data []byte) error {

	secured := "0" // 1=YES 0=NO 受信したIPパケットを構成するMACフレームを暗号化するか（原文ママ）

	cmd := fmt.Sprintf("SKSENDTO 1 %s 0E1A %s %04X", ipv6, secured, len(data)) // DATALENは16進数
	cmd = cmd + " "                                                            // DATALENのあとにスペースを入れないとデータ入力待ちにならない
	// log.Debug().Msgf("--> %s", cmd)
	n, err := port.Write([]byte(cmd))
	if err != nil {
		log.Err(err).Msgf("skSendTo: command send error")
		return err
	}
	if n != len(cmd) {
		return fmt.Errorf("skSendTo: unexpected command sent bytes. expected=%d actual=%d", len(cmd), n)
	}

	// データ入力待ちになるまでちょっと間隔を空ける
	time.Sleep(100 * time.Millisecond)

	n, err = port.Write(data)
	if err != nil {
		log.Err(err).Msgf("skSendTo: Data send error")
		return err
	}
	if n != len(data) {
		return fmt.Errorf("skSendTo: unexpected data sent bytes. expected=%d actual=%d", len(data), n)
	}

	log.Debug().Msgf("--> %s%s", cmd, hex.EncodeToString(data))

	return nil
}
