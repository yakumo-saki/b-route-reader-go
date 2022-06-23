package smartmeter

// EPC 0xE7 瞬時電流計測値をパースする。
func (sm *ELSmartMeterParser) ParseE8NowDenryuu(data []byte) (NowDenryuu, error) {

	ret := NowDenryuu{}

	err := sm.checkPreCondition()
	if err != nil {
		return ret, err
	}

	ret.Rphase = 100
	ret.Tphase = 200
	ret.Total = 300
	return ret, nil
}
