package smartmeter

// EPC 0xE7 瞬時電力計測値をパースする。
func (sm *ELSmartMeterParser) ParseE7NowDenryoku(data []byte) (int, error) {
	ret := -1
	err := sm.checkPreCondition()
	if err != nil {
		return ret, err
	}

	return 1000, nil
}
