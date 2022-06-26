package smartmeter

import (
	"bytes"
	"encoding/binary"
)

// EPC 0xE7 瞬時電力計測値をパースする。単位 W
func (sm *ELSmartMeterParser) ParseE7NowDenryoku(data []byte) (int, error) {

	var watt int32
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &watt)
	if err != nil {
		return -1, err
	}

	return int(watt), nil
}
