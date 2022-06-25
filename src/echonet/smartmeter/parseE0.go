package smartmeter

// EPC 0xE0 積算電力量をパースする。
// TODO mock
func (sm *ELSmartMeterParser) ParseE0DeltaDenryoku(data []byte) (float64, error) {
	ret := float64(-1)

	err := sm.checkPreCondition()
	if err != nil {
		return ret, err
	}

	return 1000 * sm.unit, nil
}
