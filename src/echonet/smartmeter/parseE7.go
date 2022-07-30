package smartmeter

import (
	"bytes"
	"encoding/binary"

	"github.com/shopspring/decimal"
)

// EPC 0xE7 瞬時電力計測値をパースする。単位 W
func (sm *ELSmartMeterParser) ParseE7NowDenryoku(data []byte) (decimal.Decimal, error) {

	var watt int32
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &watt)
	if err != nil {
		return ERR_VAL, err
	}

	return decimal.NewFromInt32(watt), nil
}
