package echonet

import (
	"encoding/hex"
	"fmt"

	"github.com/rs/zerolog/log"
)

// プロパティリスト要求電文を作成する。
//
func CreateEchonetInfMessage(transactionId int, properties []string) []byte {
	result := []byte{}
	result = append(result, []byte{0x10, 0x81}...)  // EHD 1081
	result = append(result, []byte{0x56, 0x78}...)  // Transaction ID
	result = append(result, DEVICE_CONTROLLER...)   // SEOJ
	result = append(result, DEVICE_METER...)        // DEOJ
	result = append(result, []byte{OPC_INF_REQ}...) // ESV: INF
	result = append(result, []byte{0x00}...)        // OPC: 0
	// result = append(result, []byte{PROPERTY_NOW_DENRYUU}...) // EPC: test
	// result = append(result, []byte{0x00}...)                 // EDT LENGTH

	return result
}

// プロパティ値要求電文を生成する
// SKSENDTOで使用できるように []byteで返す
//   EHD  1081 ... Echonet lite 固定値 2byte
// TID  1234 ... Transaction ID。応答に同じIDをつけてくれる 2byte
// SEOJ 05FF01 ... 送信元機器種別。 とりあえずコントローラ 05FF01 固定 3byte
// DEOJ 028801 ... 送信先機器種別。低圧スマートメーター 028801 3byte
// ESV         ... 命令(get,set)を入れる。 GET=62 1byte
// OPC         ... 命令数
// EPC         ... プロパティコード
// PDC         ... EDTデータ長。 GET時は 0  1byte
// EDT         ... データ部
func CreateEchonetGetMessage(transactionId uint16, properties []byte) []byte {
	tid, err := hex.DecodeString(fmt.Sprintf("%04d", transactionId))
	if err != nil {
		tid = []byte{0xEE, 0xEE}
	}

	opc, err := hex.DecodeString(fmt.Sprintf("%02x", len(properties)))
	if err != nil {
		log.Err(err).Msgf("Failed to convert %d to hex", len(properties))
	}

	result := []byte{}
	result = append(result, []byte{0x10, 0x81}...) // EHD 1081
	result = append(result, tid...)                // Transaction ID
	result = append(result, DEVICE_CONTROLLER...)  // SEOJ
	result = append(result, DEVICE_METER...)       // DEOJ
	result = append(result, []byte{OPC_GET}...)    // ESV: GET
	result = append(result, opc...)                // OPC: 1
	for _, v := range properties {
		result = append(result, []byte{v}...)    // EPC: test
		result = append(result, []byte{0x00}...) // PDC: EDT LENGTH
		// EDT: when GET it is none
	}

	return result
}
