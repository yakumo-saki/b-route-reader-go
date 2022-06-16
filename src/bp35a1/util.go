package bp35a1

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

// BP35A1にコマンドを送信する。自動的に末尾に改行（CRLF）を追加して送信する。
func sendCommand(cmd string) error {
	n, err := port.Write([]byte(fmt.Sprintf("%s\r\n", cmd)))
	if err != nil {
		command := strings.Split(cmd, " ")[0]
		log.Err(err).Msgf("Command send error: %s", command)
		return err
	}

	if n != (len(cmd) + 2) {
		command := strings.Split(cmd, " ")[0]

		fmt.Errorf("Unexpected sent bytes. expected=%d actual=%d command=%s",
			len(cmd)+2, n, command)
	}
	return nil
}

func containsOK(ret []string) bool {
	return true
}

func containsEvent22(ret []string) bool {
	return true
}
