package smartmeter

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestELSmartMeterParser_ParseAndStoreSimple(t *testing.T) {
	assert := assert.New(t)

	hexE1, _ := hex.DecodeString("04")
	hexD3, _ := hex.DecodeString("00010000")

	parser := ELSmartMeterParser{}
	parser.ParseAndStoreE1DeltaUnit(hexE1)
	parser.ParseAndStoreD3Multiplier(hexD3)
	assert.Equal(10000, parser.multiplier)
	assert.Equal(0.0001, parser.unit)
}

func TestELSmartMeterParser_E1Fail(t *testing.T) {
	assert := assert.New(t)

	lenErr, _ := hex.DecodeString("00000004")
	valueErr, _ := hex.DecodeString("05")

	parser := ELSmartMeterParser{}
	_, err := parser.ParseAndStoreE1DeltaUnit(lenErr)
	assert.NotNil(err, "must be length error")

	_, err = parser.ParseAndStoreE1DeltaUnit(valueErr)
	assert.NotNil(err, "must be value error")

}
