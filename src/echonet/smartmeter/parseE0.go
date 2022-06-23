package smartmeter

// EPC 0xE0 積算電力量をパースする。
func (sm *ELSmartMeterParser) ParseE0DeltaDenryoku(data []byte) (int, error) {
	ret := -1

	err := sm.checkPreCondition()
	if err != nil {
		return ret, err
	}

	return 1000 * sm.unit, nil
}
