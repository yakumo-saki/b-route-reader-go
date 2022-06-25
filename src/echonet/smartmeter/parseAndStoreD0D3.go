package smartmeter

import (
	"encoding/hex"
	"fmt"
	"strconv"
)

// EPC 0xD3 係数
// そのまま10進数として読めばOK 000000~999999
func (sm *ELSmartMeterParser) ParseAndStoreD3Multiplier(data []byte) (int, error) {
	dataStr := hex.EncodeToString(data) // [0x01,0x02] -> "0102"

	ret64, err := strconv.ParseInt(dataStr, 10, 16)
	if err != nil {
		return -1, err
	}
	ret := int(ret64)
	sm.multiplier = int(ret)
	return int(ret), nil
}

// EPC 0xD0 係数
// 0x00：1kWh
// 0x01：0.1kWh
// 0x02：0.01kWh
// 0x03：0.001kWh
// 0x04：0.0001kWh
// 0x0A：10kWh
// 0x0B：100kWh
// 0x0C：1000kWh
// 0x0D：10000kWh
func (sm *ELSmartMeterParser) ParseAndStoreE1DeltaUnit(data []byte) (float64, error) {
	if len(data) != 1 {
		return -1, fmt.Errorf("property E1 length != 1")
	}

	ret := -1.0
	switch data[0] {
	case 0x00:
		ret = 1
	case 0x01:
		ret = 0.1
	case 0x02:
		ret = 0.01
	case 0x03:
		ret = 0.001
	case 0x04:
		ret = 0.0001
	case 0x0A:
		ret = 10
	case 0x0B:
		ret = 100
	case 0x0C:
		ret = 1000
	case 0x0D:
		ret = 10000
	default:
		return -1, fmt.Errorf("property E1 unknown value %02d", data[0])
	}

	sm.unit = ret
	return ret, nil
}
