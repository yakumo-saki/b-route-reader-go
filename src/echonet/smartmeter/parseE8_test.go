package smartmeter

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestELSmartMeterParser_ParseE8NowDenryuu(t *testing.T) {

	t.Run("通常値", func(t *testing.T) {
		parser := ELSmartMeterParser{multiplier: decimal.NewFromInt(1), unit: decimal.RequireFromString("1.0")}
		amps, err := parser.ParseE8NowDenryuu([]byte{0x00, 0x12, 0x00, 0x34})
		assert.Equal(t, decimal.RequireFromString("1.8"), amps.Rphase, "R相")
		assert.Equal(t, decimal.RequireFromString("5.2"), amps.Tphase, "T相")
		assert.Equal(t, decimal.RequireFromString("7.0"), amps.Total, "合計")
		assert.Nil(t, err)
	})
}
