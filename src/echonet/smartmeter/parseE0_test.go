package smartmeter

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestELSmartMeterParser_ParseE0DeltaDenryoku(t *testing.T) {

	t.Run("係数1_単位1.0", func(t *testing.T) {
		parser := ELSmartMeterParser{multiplier: decimal.NewFromInt(1), unit: decimal.RequireFromString("1.0")}
		kwh, err := parser.ParseE0DeltaDenryoku([]byte{0x00, 0x00, 0x00, 0x01})
		assert.Equal(t, decimal.RequireFromString("1.0"), kwh)
		assert.Nil(t, err)
	})

	t.Run("係数2_単位1.0", func(t *testing.T) {
		parser := ELSmartMeterParser{multiplier: decimal.NewFromInt(2), unit: decimal.RequireFromString("1.0")}
		kwh, err := parser.ParseE0DeltaDenryoku([]byte{0x00, 0x00, 0x00, 0x01})
		assert.Equal(t, decimal.RequireFromString("2.0"), kwh)
		assert.Nil(t, err)
	})

	t.Run("係数1_単位0.0001", func(t *testing.T) {
		parser := ELSmartMeterParser{multiplier: decimal.NewFromInt(1), unit: decimal.RequireFromString("0.0001")}
		kwh, err := parser.ParseE0DeltaDenryoku([]byte{0x00, 0x00, 0x00, 0x01})
		assert.Equal(t, decimal.RequireFromString("0.0001"), kwh)
		assert.Nil(t, err)
	})

}
