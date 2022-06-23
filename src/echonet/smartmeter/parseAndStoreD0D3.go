package smartmeter

// EPC 0xD3 係数
func (sm *ELSmartMeterParser) ParseAndStoreD3Multiplier(data []byte) (int, error) {
	ret := 1

	sm.multiplier = ret
	return ret, nil
}

// EPC 0xD0 係数
func (sm *ELSmartMeterParser) ParseAndStoreD0DeltaUnit(data []byte) (int, error) {
	ret := 1

	sm.unit = ret
	return ret, nil
}
