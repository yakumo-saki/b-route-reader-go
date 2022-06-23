package echonet

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

// 1081〜から始まるASCII表記のEchonetLite電文を解釈する
func Parse(echonetLiteMessage string) (EchonetLite, error) {

	el := NewEchonetLite()

	if !strings.HasPrefix(echonetLiteMessage, "1081") {
		return el, fmt.Errorf("unknown echonet lite message (prefix not 1081) %s", echonetLiteMessage)
	}

	// esv

	el.Ehd = echonetLiteMessage[0:4]
	el.TransactionId = echonetLiteMessage[4:8]
	el.Seoj = echonetLiteMessage[8:14]
	el.Deoj = echonetLiteMessage[14:20]
	esv, err := stringToSingleByte(echonetLiteMessage[20:22])
	if err != nil {
		return el, err
	}
	el.Esv = esv
	opc, err := strconv.ParseInt(echonetLiteMessage[22:24], 16, 16)
	if err != nil {
		return el, err
	}
	el.Opc = int(opc)

	epcStart := 24
	for i := 0; i < el.Opc; i++ {
		data := echonetLiteOneData{}

		// EPC (プロパティ)
		epc, err := substring(echonetLiteMessage, epcStart, 2)
		if err != nil {
			return el, err
		}
		data.Epc, err = stringToSingleByte(epc)
		if err != nil {
			return el, err
		}

		// PDC（EDTバイト数）
		pdcStart := epcStart + 2
		pdcStr, err := substring(echonetLiteMessage, pdcStart, 2)
		if err != nil {
			return el, err
		}
		pdc, err := strconv.ParseInt(pdcStr, 16, 16)
		if err != nil {
			return el, err
		}
		data.Pdc = int(pdc)

		// EDT
		edtStart := pdcStart + 2
		edt, err := substring(echonetLiteMessage, edtStart, data.Pdc*2) // pdcはバイト数。ASCII表記だと2倍
		if err != nil {
			return el, err
		}
		data.Edt, err = hex.DecodeString(edt)
		if err != nil {
			return el, err
		}

		epcStart = edtStart + (data.Pdc * 2)

		el.Properties[data.Epc] = data.Edt
	}

	return el, nil
}

func stringToSingleByte(hexString string) (byte, error) {
	if len(hexString) != 2 {
		return byte(0x00), fmt.Errorf("string length is not 2")
	}
	by, err := hex.DecodeString(hexString)
	if err != nil {
		return byte(0x00), err
	}
	if len(by) != 1 {
		return byte(0x00), fmt.Errorf("BUG: byte length != 1")
	}

	return by[0], nil
}

func substring(str string, start, length int) (string, error) {
	if len(str) < (start + length) {
		return "", fmt.Errorf("insufficient string length. need=%d actual=%d",
			(start + length), len(str))
	}

	ret := []rune(str)[start:(start + length)] // unicode対応。不要だけども

	return string(ret), nil
}
