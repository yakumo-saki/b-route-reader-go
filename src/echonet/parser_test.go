package echonet_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakumo-saki/b-route-reader-go/src/echonet"
)

func TestParserGetSingleProp(t *testing.T) {
	assert := assert.New(t)
	elmsg := "1081232802880105FF017201D30400000001"
	el, err := echonet.Parse(elmsg)

	assert.Nil(err, "err != nil")
	assert.Equal("2328", el.TransactionId, "TID != 2328")
	assert.Equal("028801", el.Seoj, "SEOJ != 028801")
	assert.Equal("05FF01", el.Deoj, "DEOJ != 05FF01")
	assert.Equal(byte(0x72), el.Esv, "ESV != 0x72")
	assert.Equal(1, el.Opc, "OPC != 1")
	assert.Equal(1, len(el.Properties), "el.Properties != el.Opc")

	dat := el.Properties[0]
	assert.Equal(byte(0xD3), dat.Epc, "EPC1 != 0xD3")
	assert.Equal(4, dat.Pdc, "PDC1 != 4")
	assert.Equal([]byte{0x00, 0x00, 0x00, 0x01}, dat.Edt, "EDT1 != 0x00000001")

}

func TestParserGetMultiProp(t *testing.T) {
	assert := assert.New(t)
	elmsg := "1081A00102880105FF017302EA0B07E606150F1E0000016CA7EB0B07E606150F1E0000000005"
	el, err := echonet.Parse(elmsg)

	assert.Nil(err, "err != nil")
	assert.Equal("A001", el.TransactionId, "TID != A001")
	assert.Equal("028801", el.Seoj, "SEOJ != 028801")
	assert.Equal("05FF01", el.Deoj, "DEOJ != 05FF01")
	assert.Equal(byte(0x73), el.Esv, "ESV != 0x72")
	assert.Equal(el.Opc, 2, "OPC != 2")
	assert.Equal(el.Opc, len(el.Properties), "el.Properties != el.Opc")

	dat, ok := el.Properties[byte(0xEA)]
	assert.True(ok)
	hx, _ := hex.DecodeString("07E606150F1E0000016CA7")
	assert.Equal(dat, hx, "EDT EA != 0x07E606150F1E0000016CA7")

	dat, ok := el.Properties[byte(0xEA)]
	assert.True(ok)
	hx, _ := hex.DecodeString("07E606150F1E0000016CA7")
	assert.Equal(dat, hx, "EDT EA != 0x07E606150F1E0000016CA7")

	dat, ok := el.Properties[byte(0xEB)]
	hx, _ = hex.DecodeString("07E606150F1E0000000005")
	assert.Equal(dat, hx, "EDT EB != 07E606150F1E0000000005")

}
