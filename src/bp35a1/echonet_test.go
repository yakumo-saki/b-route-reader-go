package bp35a1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_findFirstEchonetUnicast(t *testing.T) {
	assert := assert.New(t)

	testcase := []string{
		"EVENT 21 FE80:0000:0000:0000:021C:6400:03CD:76A4 00",
		"OK",
		"ERXUDP FE80:0000:0000:0000:0000:0000:5678:1234 FE80:0000:0000:0000:1234:5678:ABCD:12EF 0E1A 0E1A 000000000000001234 1 0018 1081100202880105FF017202E7040000029CE80400280028",
	}

	ret := findEchonetResponses(testcase, "1002")

	assert.Equal(1, len(ret))
	assert.Equal(testcase[2], ret[0])
}

func Test_isEchonetUnicastResponse(t *testing.T) {
	assert := assert.New(t)

	testcase := []string{
		"EVENT 21 FE80:0000:0000:0000:021C:6400:03CD:76A4 00",
		"OK",
		"ERXUDP FE80:0000:0000:0000:0000:0000:5678:1234 FE80:0000:0000:0000:1234:5678:ABCD:12EF 0E1A 0E1A 000000000000001234 1 0018 1081100202880105FF017202E7040000029CE80400280028",
	}

	tid := "1002"
	assert.False(isEchonetUnicastResponse(testcase[0], tid))
	assert.False(isEchonetUnicastResponse(testcase[1], tid))
	assert.True(isEchonetUnicastResponse(testcase[2], tid))
}
