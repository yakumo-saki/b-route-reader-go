package smartmeter

import (
	"bytes"
	"encoding/binary"
)

// EPC 0xE0 積算電力量をパースする。
// 積算電力量を扱うので、係数(D3)と単位(E1)が必要
func (sm *ELSmartMeterParser) ParseE0DeltaDenryoku(data []byte) (float64, error) {

	err := sm.checkPreCondition()
	if err != nil {
		return -1, err
	}

	// numStr := hex.EncodeToString(data)
	// num, err := strconv.ParseFloat(numStr, 32)
	var val uint32
	err = binary.Read(bytes.NewBuffer(data), binary.BigEndian, &val)
	if err != nil {
		return -1, err
	}

	kwh := float64(val) * float64(sm.multiplier) * sm.unit

	return kwh, nil
}
