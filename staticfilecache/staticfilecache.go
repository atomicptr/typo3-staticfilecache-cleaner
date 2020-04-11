package staticfilecache

import (
	"encoding/json"
	"log"
	"strings"
	"time"
)

// CacheEntry represents the content of a cache entry file like "index.config.json"
type CacheEntry struct {
	Generated string              `json:"generated"`
	Headers   map[string][]string `json:"headers"`
}

// Parse parses the data within a cache entry file and returns it
func Parse(data []byte) (*CacheEntry, error) {
	var entry CacheEntry
	err := json.Unmarshal(data, &entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

// GetHeader returns the header value associated with a given key
func (c *CacheEntry) GetHeader(headerName string) string {
	headerValues, ok := c.Headers[headerName]
	if !ok {
		return ""
	}
	return strings.Join(headerValues, ",")
}

func (c *CacheEntry) IsExpired() bool {
	expiresValue := c.GetHeader("Expires")

	// no expires header set? Assume something is wrong and mark this as expired
	if expiresValue == "" {
		return true
	}

	date, err := time.Parse(time.RFC1123, expiresValue)
	// couldn't parse date? Assume something is wrong and mark this as expired
	if err != nil {
		log.Println(err)
		return true
	}

	return time.Now().After(date)
}
