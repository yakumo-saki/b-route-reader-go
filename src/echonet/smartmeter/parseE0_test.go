package smartmeter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestELSmartMeterParser_ParseE0DeltaDenryoku(t *testing.T) {

	t.Run("係数1_単位1.0", func(t *testing.T) {
		parser := ELSmartMeterParser{multiplier: 1, unit: 1.0}
		kwh, err := parser.ParseE0DeltaDenryoku([]byte{0x00, 0x00, 0x00, 0x01})
		assert.Equal(t, 1.0, kwh)
		assert.Nil(t, err)
	})

	t.Run("係数2_単位1.0", func(t *testing.T) {
		parser := ELSmartMeterParser{multiplier: 2, unit: 1.0}
		kwh, err := parser.ParseE0DeltaDenryoku([]byte{0x00, 0x00, 0x00, 0x01})
		assert.Equal(t, 1.0*2, kwh)
		assert.Nil(t, err)
	})

	t.Run("係数1_単位0.0001", func(t *testing.T) {
		parser := ELSmartMeterParser{multiplier: 1, unit: 0.0001}
		kwh, err := parser.ParseE0DeltaDenryoku([]byte{0x00, 0x00, 0x00, 0x01})
		assert.Equal(t, 1.0*1*0.0001, kwh)
		assert.Nil(t, err)
	})

}
