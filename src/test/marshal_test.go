package test

import (
	"encoding/json"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/yakumo-saki/b-route-reader-go/src/bp35a1"
)

func TestMarshal(t *testing.T) {
	assert := assert.New(t)

	data := bp35a1.ElectricData{}
	data["E0"] = decimal.RequireFromString("10190.40005")

	jsonMap := map[string]interface{}{}
	for k, v := range data {
		jsonMap[k] = v
	}

	json, err := json.Marshal(jsonMap)
	jsonStr := string(json)

	assert.Nil(err)
	assert.Equal("{\"E0\":\"10190.40005\"}", jsonStr)

}
