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

var P_DELTA_DENRYOKU = byte(0xE0)   // E0 積算電力量計測値（正方向）
var P_DELTA_DENRYOKU_R = byte(0xE3) // E3 積算電力量計測値（逆方向）
var P_NOW_DENRYOKU = byte(0xE7)     // E7瞬時電力。 4byte HEX signed long
var P_NOW_DENRYUU = byte(0xE8)      // E8瞬時電流。2byte HEX * 2(R相 T相) 0.1A単位。T相が0x7FFDの場合単相2線式
var P_DELTA_HISTORY = byte(0xE2)

// 係数。D3 すべての値にこの値を乗算する必要がある。応答は10進数
var P_MULTIPLIER = byte(0xD3)

// 積算電力量単位。 E1 1byte hex. 00=1kWh 01=0.1kWh 04=0.0001kWh 0A=10kWh 0B=100kWh 0C=1000kWh 0D=10000kWh
var P_DELTA_UNIT = byte(0xE1)
