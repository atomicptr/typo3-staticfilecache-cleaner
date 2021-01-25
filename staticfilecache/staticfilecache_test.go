package staticfilecache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseInvalid(t *testing.T) {
	jsonString := `{"test": asdf`

	_, err := Parse([]byte(jsonString))
	assert.NotNil(t, err)
}

const expiredJSONString = `{
	"generated": "Sat, 11 Apr 2020 14:28:05 CEST",
	"headers": {
		"Single-Value-Header": [
			"Value"
		],
		"Multi-Value-Header": [
			"Value 1",
			"Value 2"
		],
		"Expires": [
			"Sat, 11 Apr 2020 14:28:05 CEST"
		]
	}
}`

func TestParse(t *testing.T) {
	_, err := Parse([]byte(expiredJSONString))
	assert.Nil(t, err)
}

func TestGetHeader(t *testing.T) {
	cacheEntry, err := Parse([]byte(expiredJSONString))
	assert.Nil(t, err)

	// invalid header
	assert.Equal(t, "", cacheEntry.GetHeader("Invalid-Test-Header"))

	// valid headers
	assert.Equal(t, "Value", cacheEntry.GetHeader("Single-Value-Header"))
	assert.Equal(t, "Value 1,Value 2", cacheEntry.GetHeader("Multi-Value-Header"))
	assert.Equal(t, "Sat, 11 Apr 2020 14:28:05 CEST", cacheEntry.GetHeader("Expires"))
}

func TestIsExpired(t *testing.T) {
	cacheEntry, err := Parse([]byte(expiredJSONString))
	assert.Nil(t, err)
	assert.True(t, cacheEntry.IsExpired())
}

const notExpiredJSONString = `{
	"generated": "Sat, 11 Apr 2020 14:28:05 CEST",
	"headers": {
		"Expires": [
			"Mon, 18 Mar 2120 14:28:05 CEST"
		]
	}
}`

func TestIsExpiredNotExpired(t *testing.T) {
	cacheEntry, err := Parse([]byte(notExpiredJSONString))
	assert.Nil(t, err)
	assert.False(t, cacheEntry.IsExpired())
}

const invalidDateJSONString = `{
	"generated": "Sat, 11 Apr 2020 14:28:05 CEST",
	"headers": {
		"Expires": [
			"This is not a valid date!"
		]
	}
}`

func TestIsExpiredInvalidDateString(t *testing.T) {
	cacheEntry, err := Parse([]byte(invalidDateJSONString))
	assert.Nil(t, err)
	assert.True(t, cacheEntry.IsExpired())
}

const noHeadersJSONString = `{
	"generated": "Sat, 11 Apr 2020 14:28:05 CEST",
	"headers": {}
}`

func TestIsExpiredNoExpiresHeader(t *testing.T) {
	cacheEntry, err := Parse([]byte(noHeadersJSONString))
	assert.Nil(t, err)
	assert.True(t, cacheEntry.IsExpired())
}
