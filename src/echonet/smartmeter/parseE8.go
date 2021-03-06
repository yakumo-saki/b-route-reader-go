package smartmeter

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/shopspring/decimal"
)

// EPC 0xE7 瞬時電流計測値をパースする
func (sm *ELSmartMeterParser) ParseE8NowDenryuu(data []byte) (NowDenryuu, error) {

	ret := NowDenryuu{}

	if len(data) != 4 {
		return ret, fmt.Errorf("data is not 4bytes")
	}

	rPhase, err := parseSignedShort(data[0:2])
	if err != nil {
		fmt.Println(err)
		return ret, fmt.Errorf("failed to parse Rphase value 0x%02X%02X: %w", data[0], data[1], err)
	}
	tPhase, err := parseSignedShort(data[2:4])
	if err != nil {
		fmt.Println(err)
		return ret, fmt.Errorf("failed to parse Tphase value 0x%02X%02X: %w", data[2], data[3], err)
	}

	keisuu, _ := decimal.NewFromString("0.1")
	ret.Rphase = decimal.NewFromInt(int64(rPhase)).Mul(keisuu)
	ret.Tphase = decimal.NewFromInt(int64(tPhase)).Mul(keisuu)
	ret.Total = ret.Rphase.Add(ret.Tphase)
	return ret, nil
}

func parseSignedShort(bin []byte) (int16, error) {
	var ret int16
	err := binary.Read(bytes.NewReader(bin), binary.BigEndian, &ret)
	if err != nil {
		return -1, err
	}
	return ret, nil
}
