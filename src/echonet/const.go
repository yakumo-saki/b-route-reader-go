package echonet

// DEOJに指定する低圧スマートメーターのデバイスID(028801)
var DEVICE_METER = []byte{0x02, 0x88, 0x01}

// SEOJに指定するコントローラーのデバイスID(05FF01)
var DEVICE_CONTROLLER = []byte{0x05, 0xFF, 0x01}

// ノードプロファイルのデバイスID(0EF001)
// これはデバイスではなく、INFコマンドの応答で使われるデバイスID。
var DEVICE_NODE_PROFILE = []byte{0x0E, 0xF0, 0x01}

// プロパティ取得命令
var OPC_GET = byte(0x62)     // プロパティ取得
var OPC_INF_REQ = byte(0x63) // インスタンスリスト通知要求

var OPC_INF = byte(0x73) // インスタンスリスト通知

//
var P_GET_PROPERTY_MAP = byte(0x9F)
var P_SERIAL_NO = byte(0x8D)

var P_DELTA_DENRYOKU = byte(0xE0)
var P_NOW_DENRYOKU = byte(0xE7)
var P_NOW_DENRYUU = byte(0xE8)
var P_DELTA_HISTORY = byte(0xE2)
var P_DELTA_UNIT = byte(0xE1)

var P_UNIT = byte(0xD3)
