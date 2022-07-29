package bp35a1

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

var STD_TIMEOUT_DURATION = 15 * time.Second
var LONG_TIMEOUT_DURATION = 60 * time.Second

//
func waitForResultStrings(stopWords []string, timeoutDuration time.Duration) ([]string, error) {
	stopFn := func(line string) bool {
		if len(stopWords) == 0 {
			return true
		}
		for _, sw := range stopWords {
			if strings.Contains(line, sw) {
				return true
			}
		}
		return false
	}
	return waitForResult(stopFn, timeoutDuration)
}

// OK等を返すコマンドの応答を返す
func waitForResultOK() ([]string, error) {
	return waitForResultStrings(RET_STOP_WORDS, STD_TIMEOUT_DURATION)
}

// SKLL64の応答を返す
// このコマンドはいきなりIPv6アドレスだけを返してくる
func waitForResultSKLL64() ([]string, error) {
	return waitForResultStrings([]string{}, STD_TIMEOUT_DURATION)
}

func waitForResultSKSCAN() ([]string, error) {
	return waitForResultStrings([]string{RET_SCAN_COMPLETE}, LONG_TIMEOUT_DURATION)
}

func waitForResultSKJOIN() ([]string, error) {
	return waitForResultStrings([]string{RET_JOIN_COMPLETE}, LONG_TIMEOUT_DURATION)
}

func waitForResultERXUDP(tid string) ([]string, error) {
	stopFn := func(line string) bool {
		return isEchonetUnicastResponse(line, tid)
	}
	return waitForResult(stopFn, LONG_TIMEOUT_DURATION)
}

func waitForResult(completeCheckFunc func(string) bool, timeoutDuration time.Duration) ([]string, error) {

	log.Debug().Msgf("Response start. timeout=%s", timeoutDuration)
	BYTE_CR := []byte("\r")
	BYTE_LF := []byte("\n")
	BYTE_CRLF := []byte("\r\n")
	LF := byte(0xa)

	port.SetReadTimeout(300 * time.Millisecond)

	timeoutTime := time.Now().Add(timeoutDuration)

	var result []string

	// 応答を全部溜め込むバッファ
	var byteBuf []byte

	stopFlag := false

	for {
		// 一回のreadで読んだデータのバッファ
		readBuf := make([]byte, 200)

		n, err := port.Read(readBuf)
		if err != nil {
			log.Err(err).Msg("Read error")
			return nil, err
		}
		if n == 0 && len(byteBuf) == 0 && stopFlag {
			if len(result) > 0 {
				break
			}
		}

		// タイムアウトまわり
		if n == 0 {
			if time.Now().After(timeoutTime) {
				return result, fmt.Errorf("waitForResult timeout reached")
			}
		} else {
			timeoutTime = time.Now().Add(timeoutDuration) // タイムアウト延長
		}

		byteBuf = append(byteBuf, readBuf[:n]...)

		// byteBufを改行で分割する
		for {
			if bytes.Contains(byteBuf, BYTE_CR) {
				crPos := bytes.Index(byteBuf, BYTE_CR)
				lineBuf := byteBuf[:crPos]
				lineStr := trimResponse(string(lineBuf))

				byteBuf = byteBuf[crPos:] // 改行コードは削除

				// CRLFで区切られている場合、LFが残るので削除
				switch {
				case bytes.HasPrefix(byteBuf, BYTE_CRLF):
					byteBuf = byteBuf[len(BYTE_CRLF):]
					log.Debug().Msgf("<-- %s<CRLF>", lineStr)
				case bytes.HasPrefix(byteBuf, BYTE_LF):
					byteBuf = byteBuf[len(BYTE_LF):]
					log.Debug().Msgf("<-- %s<LF>", lineStr)
				case bytes.HasPrefix(byteBuf, BYTE_CR):
					byteBuf = byteBuf[len(BYTE_CR):]
					log.Debug().Msgf("<-- %s<CR>", lineStr)
				}

				result = append(result, lineStr)

				stopFlag = (stopFlag || completeCheckFunc(lineStr))

			} else if len(byteBuf) == 1 && byteBuf[0] == LF {
				// SKLL64の時、LFだけがバッファに残ってしまい無限ループすることへの対策
				// 本来は起きないはずなのだが…
				// log.Warn().Msgf("executing start with LF workaround")
				byteBuf = byteBuf[1:]
				break
			} else {
				// 改行コードが含まれない
				// 行の途中でバッファがいっぱいになったか、応答の末尾まで読み切った
				// どちらにしてももう一度readを呼ぶ
				break
			}

		} // 無限ループ
	}

	log.Debug().Msg("Response done")
	return result, nil
}

func trimResponse(str string) string {
	result := str
	result = strings.TrimPrefix(result, " ")
	result = strings.TrimPrefix(result, "\r")
	result = strings.TrimPrefix(result, "\n")

	return result

}
