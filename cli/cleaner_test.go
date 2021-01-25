package cli

import (
	"fmt"
	"github.com/atomicptr/typo3-staticfilecache-cleaner/staticfilecache"
	copylib "github.com/otiai10/copy"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
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

func TestIntegrationCleanPathInDryRun(t *testing.T) {
	const numFilesDeleted = 5

	flagDryRun = true

	deleted := cleanPath(testDataBaseDir)
	assert.Equal(t, numFilesDeleted, deleted)
}

func TestIntegrationCleanPath(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "typo3-staticfilecache-cleaner-testdir")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	cwd, err := os.Getwd()
	assert.NoError(t, err)

	fixtureDir := filepath.Join(cwd, "..", "testdata", "typo3temp")

	err = copylib.Copy(fixtureDir, dir)
	assert.NoError(t, err)

	// assert that service directory actually exists
	serviceDir := filepath.Join(dir, "tx_staticfilecache", "http_domain.com_80", "service")
	assert.DirExists(t, serviceDir)

	flagDryRun = false

	cleanPaths([]string{dir})

	// test if directory tx_staticfilecache/service was also deleted
	_, err = os.Lstat(serviceDir)
	assert.True(t, os.IsNotExist(err), fmt.Sprintf("directory '%s' should not exist", serviceDir))
}
