package smartmeter

import "fmt"

type ELSmartMeterParser struct {
	// 0xD3 係数。積算電力量計測値、履歴にこの係数をかける
	multiplier int

	// 0xE1 積算電力量単位 積算電力量にこの値を掛けると kWh になる。
	// 0.0001 〜 10000 までの値を取る(10^nのみで中途半端な値はない)
	unit int
}

func (sm *ELSmartMeterParser) checkPreCondition() error {
	if sm.multiplier == 0 {
		return fmt.Errorf("multiplier (0xD3) is not set. parseE3 first")
	}
	if sm.unit == 0 {
		return fmt.Errorf("unit (0xD1) is not set. parseE1 first")
	}

	return nil
}
