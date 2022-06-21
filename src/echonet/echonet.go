package echonet

// プロパティリスト要求電文を作成する。
//
func CreateEchonetInfMessage(transactionId int, properties []string) []byte {
	result := []byte{}
	result = append(result, []byte{0x10, 0x81}...) // EHD 1081
	result = append(result, []byte{0x13, 0x57}...) // Transaction ID
	result = append(result, DEVICE_CONTROLLER...)  // SEOJ
	result = append(result, DEVICE_METER...)       // DEOJ
	result = append(result, []byte{OPC_INF}...)    // ESV: INF
	// result = append(result, []byte{0x00}...)                 // OPC: 0
	// result = append(result, []byte{PROPERTY_NOW_DENRYUU}...) // EPC: test
	// result = append(result, []byte{0x00}...)                 // EDT LENGTH

	return result
}

// プロパティ値要求電文を生成する
// SKSENDTOで使用できるように []byteで返す
// EHD  1081 ... Echonet lite 固定値 2byte
// TID  1234 ... Transaction ID。応答に同じIDをつけてくれる 2byte
// SEOJ 05FF01 ... 送信元機器種別。 とりあえずコントローラ 05FF01 固定 3byte
// DEOJ 028801 ... 送信先機器種別。低圧スマートメーター 028801 3byte
// ESV         ... 命令(get,set)を入れる。 GET=62 1byte
// OPC         ... 命令数
// EPC         ... プロパティコード
// PDC         ... EDTデータ長。 GET時は 0  1byte
// EDT         ... データ部
func CreateEchonetGetMessage(transactionId int, properties []string) []byte {
	result := []byte{}
	result = append(result, []byte{0x10, 0x81}...)           // EHD 1081
	result = append(result, []byte{0x98, 0x76}...)           // Transaction ID
	result = append(result, DEVICE_CONTROLLER...)            // SEOJ
	result = append(result, DEVICE_METER...)                 // DEOJ
	result = append(result, []byte{OPC_GET}...)              // ESV: GET
	result = append(result, []byte{0x01}...)                 // OPC: 1
	result = append(result, []byte{PROPERTY_NOW_DENRYUU}...) // EPC: test
	result = append(result, []byte{0x00}...)                 // EDT LENGTH

	return result
}
