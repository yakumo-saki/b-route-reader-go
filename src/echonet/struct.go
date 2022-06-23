package echonet

// EchonetLite(=EL)電文
type EchonetLite struct {
	Ehd           string          // 1081固定
	TransactionId string          // トランザクションID。識別用ID
	Seoj          string          // 送信元ELオブジェクト
	Deoj          string          // 宛先ELオブジェクト
	Esv           byte            // Echonetサービスコード。GetとかGet_Resとか。
	Opc           int             // Propertiesの数
	Properties    map[byte][]byte // EPC(プロパティ) -> EDT(値)のマップ
}

func NewEchonetLite() EchonetLite {
	el := EchonetLite{}
	el.Properties = map[byte][]byte{}
	return el
}

// EchonetLite電文 プロパティ繰り返し部
// 解釈時に一時的に使用する。 EL.Properties[Epc] = Edt するだけ。
// PDCは解釈時には有用だが、結局の所 len(Edt) なので解釈後は不要なので捨てる
type echonetLiteOneData struct {
	// ELプロパティ値のキー
	Epc byte
	// EDTのバイト数
	Pdc int
	// ELプロパティ値のデータ
	Edt []byte
}
