package bp35a1

import (
	"bytes"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// OK等を返すコマンドの応答を返す
func waitForResult() ([]string, error) {
	return waitForResultImpl(RET_STOP_WORDS)
}

// SKLL64の応答を返す
// このコマンドはいきなりIPv6アドレスだけを返してくる
func waitForResultSKLL64() ([]string, error) {
	return waitForResultImpl([]string{})
}

func waitForResultSKSCAN() ([]string, error) {
	return waitForResultImpl([]string{RET_SCAN_COMPLETE})
}

func waitForResultImpl(stopWords []string) ([]string, error) {

	log.Debug().Msg("Response start")
	BYTE_CR := []byte("\r")
	BYTE_LF := []byte("\n")

	port.SetReadTimeout(300 * time.Millisecond)

	var result []string

	// 応答を全部溜め込むバッファ
	var byteBuf []byte

	// 一回のreadで読んだデータのバッファ
	readBuf := make([]byte, 100)

	stopFlag := false
	if len(stopWords) == 0 {
		stopFlag = true
	}

	for {
		// Reads up to 100 bytes
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

		byteBuf = append(byteBuf, readBuf[:n]...)

		// byteBufを改行で分割する
		for {
			if bytes.Contains(byteBuf, BYTE_CR) {
				crPos := bytes.Index(byteBuf, BYTE_CR)
				lineBuf := byteBuf[:crPos]
				lineStr := trimResponse(string(lineBuf))

				byteBuf = byteBuf[crPos+len(BYTE_CR):] // 改行コードは削除

				// CRLFで区切られている場合、LFが残るので削除
				if bytes.HasPrefix(byteBuf, BYTE_LF) {
					byteBuf = byteBuf[1:]
					log.Debug().Msgf("<-- %s<CRLF>", lineStr)
				} else {
					log.Debug().Msgf("<-- %s<CR>", lineStr)
				}

				result = append(result, lineStr)

				// ストップワード（通常、コマンド応答の末尾に来るワード）を見つけたら終了フラグを立てる
				// OK を返したあとにさらに応答を返すコマンドがあるため（しかし、そのコマンドは使わない）
				for _, stopWord := range stopWords {
					if strings.Contains(lineStr, stopWord) {
						stopFlag = true
						break
					}
				}
			} else {
				// 改行コードが含まれない
				// 行の途中でバッファがいっぱいになったか、応答の末尾まで読み切った
				// どちらにしてももう一度readを呼ぶ
				break
			}
		}
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
