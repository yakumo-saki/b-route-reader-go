package smartmeter

import (
	"bytes"
	"encoding/binary"

	"github.com/shopspring/decimal"
)

// EPC 0xE0 積算電力量をパースする。
// 積算電力量を扱うので、係数(D3)と単位(E1)が必要
func (sm *ELSmartMeterParser) ParseE0DeltaDenryoku(data []byte) (decimal.Decimal, error) {

	err := sm.checkPreCondition()
	if err != nil {
		return ERR_VAL, err
	}

	// numStr := hex.EncodeToString(data)
	// num, err := strconv.ParseFloat(numStr, 32)
	var val uint32
	err = binary.Read(bytes.NewBuffer(data), binary.BigEndian, &val)
	if err != nil {
		return ERR_VAL, err
	}

	kwh := decimal.NewFromInt(int64(val)).Mul(sm.multiplier).Mul(sm.unit)

	return kwh, nil
}
