package smartmeter

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/shopspring/decimal"
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
	sm.multiplier = decimal.NewFromInt(ret64)
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
func (sm *ELSmartMeterParser) ParseAndStoreE1DeltaUnit(data []byte) (string, error) {
	if len(data) != 1 {
		return "-1", fmt.Errorf("property E1 length != 1")
	}

	var ret = decimal.Zero
	switch data[0] {
	case 0x00:
		ret = decimal.RequireFromString("1")
	case 0x01:
		ret = decimal.RequireFromString("0.1")
	case 0x02:
		ret = decimal.RequireFromString("0.01")
	case 0x03:
		ret = decimal.RequireFromString("0.001")
	case 0x04:
		ret = decimal.RequireFromString("0.0001")
	case 0x0A:
		ret = decimal.RequireFromString("10")
	case 0x0B:
		ret = decimal.RequireFromString("100")
	case 0x0C:
		ret = decimal.RequireFromString("1000")
	case 0x0D:
		ret = decimal.RequireFromString("10000")
	default:
		return "-1", fmt.Errorf("property E1 unknown value %02d", data[0])
	}

	sm.unit = ret
	return ret.String(), nil
}
