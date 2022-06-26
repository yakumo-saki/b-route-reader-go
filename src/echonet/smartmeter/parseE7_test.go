package smartmeter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestELSmartMeterParser_ParseE7NowDenryoku(t *testing.T) {
	parser := ELSmartMeterParser{multiplier: 1, unit: 1.0}

	t.Run("最大値", func(t *testing.T) {
		watt, err := parser.ParseE7NowDenryoku([]byte{0x7F, 0xFF, 0xFF, 0xFD})
		assert.Equal(t, 2147483645, watt)
		assert.Nil(t, err)
	})

	t.Run("マイナス値", func(t *testing.T) {
		watt, err := parser.ParseE7NowDenryoku([]byte{0x80, 0x00, 0x00, 0x01})
		assert.Equal(t, -2147483647, watt)
		assert.Nil(t, err)
	})
}
