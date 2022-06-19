package bp35a1

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

// BP35A1にコマンドを送信する。自動的に末尾に改行（CRLF）を追加して送信する。
func sendCommand(cmd string) error {
	log.Debug().Msgf("Sending %s", cmd)

	n, err := port.Write([]byte(fmt.Sprintf("%s\r\n", cmd)))
	if err != nil {
		command := strings.Split(cmd, " ")[0]
		log.Err(err).Msgf("Command send error: %s", command)
		return err
	}

	if n != (len(cmd) + 2) {
		command := strings.Split(cmd, " ")[0]

		return fmt.Errorf("unexpected sent bytes. expected=%d actual=%d command=%s",
			len(cmd)+2, n, command)
	}
	return nil
}

func containsInResult(ret []string, find string) bool {
	for _, v := range ret {
		if strings.Contains(v, find) {
			return true
		}
	}

	return false
}

func endWithResult(ret []string, find string) bool {
	return strings.Contains(ret[len(ret)-1], find)
}

// 応答がOKかFAILのコマンド用のユーティリティメソッド
func waitForOKResult() error {
	ret, err := waitForResult()
	if err != nil {
		return err
	}
	if endWithResult(ret, RET_OK) {
		return nil
	}

	return fmt.Errorf("response is not %s", RET_OK)

}
