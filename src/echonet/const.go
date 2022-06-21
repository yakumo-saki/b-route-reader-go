package echonet

// DEOJに指定する低圧スマートメーターのデバイスID(028801)
var DEVICE_METER = []byte{0x02, 0x88, 0x01}

// SEOJに指定するコントローラーのデバイスID(05FF01)
var DEVICE_CONTROLLER = []byte{0x05, 0xFF, 0x01}

// ノードプロファイルのデバイスID(0EF001)
// これはデバイスではなく、INFコマンドの応答で使われるデバイスID。
var DEVICE_NODE_PROFILE = []byte{0x0E, 0xF0, 0x01}

// プロパティ取得命令
var OPC_GET = byte(0x62)

var OPC_INF = byte(0x73)

var PROPERTY_DELTA_DENRYOKU = byte(0xE0)
var PROPERTY_NOW_DENRYOKU = byte(0xE7)
var PROPERTY_NOW_DENRYUU = byte(0xE8)
var PROPERTY_DELTA_HISTORY = byte(0xE2)
var PROPERTY_DELTA_UNIT = byte(0xE1)
