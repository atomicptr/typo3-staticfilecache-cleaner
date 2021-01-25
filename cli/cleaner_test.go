package cli

import (
	"github.com/atomicptr/typo3-staticfilecache-cleaner/staticfilecache"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"testing"
)

const (
	testDataBaseDir      = "../testdata"
	testDataIndexPageDir = "../testdata/typo3temp/tx_staticfilecache/http_domain.com_80"
)

func TestIntegrationCollectCacheEntryFilesInPath(t *testing.T) {
	const numCacheEntryFilesInTestData = 3

	files := collectCacheEntryFilesInPath(testDataBaseDir)
	assert.Len(t, files, numCacheEntryFilesInTestData)
}

func TestIntegrationFindAdjacentFiles(t *testing.T) {
	const numAdjacentFilesInIndex = 3

	files := findAdjacentFiles(filepath.Join(testDataIndexPageDir, "index.config.json"))
	assert.Len(t, files, numAdjacentFilesInIndex)
	assert.Contains(t, files, filepath.Join(testDataIndexPageDir, "random-file.txt"))
}

func TestIntegrationDeleteCacheEntry(t *testing.T) {
	const numFilesDeleted = 5

	deleted := 0

	files := collectCacheEntryFilesInPath(testDataBaseDir)

	flagDryRun = true
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		assert.Nil(t, err)

		cacheEntry, err := staticfilecache.Parse(data)
		assert.Nil(t, err)

		if cacheEntry.IsExpired() {
			deleted += deleteCacheEntry(file)
		}
	}

	assert.Equal(t, numFilesDeleted, deleted)
}

func TestIntegrationCleanPath(t *testing.T) {
	const numFilesDeleted = 5

	flagDryRun = true

	deleted := cleanPath(testDataBaseDir)
	assert.Equal(t, numFilesDeleted, deleted)
}
