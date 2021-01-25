// Package cli implements a tool that removes expired staticfilecache entries
package cli

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/atomicptr/typo3-staticfilecache-cleaner/staticfilecache"
)

const cacheEntryFileName = "index.config.json"

func cleanPaths(paths []string) {
	numDeletedFiles := 0

	log.Printf("Deleting files in %v...\n", paths)

	for _, path := range paths {
		numDeletedFiles += cleanPath(path)
	}

	deleteMessage := "Deleted"

	if flagDryRun {
		deleteMessage = "(Dry Run) Would have deleted"
	}

	log.Printf("Done! %s %d files.\n", deleteMessage, numDeletedFiles)
}

func cleanPath(path string) int {
	cacheEntryFiles := collectCacheEntryFilesInPath(path)

	numDeletedFiles := 0

	for _, filePath := range cacheEntryFiles {
		data, err := ioutil.ReadFile(filepath.Clean(filePath))
		if err != nil {
			log.Printf(`error with file "%s": %s, ignoring file...`+"\n", filePath, err)
			continue
		}

		cacheEntry, err := staticfilecache.Parse(data)
		if err != nil {
			log.Printf(`error with file "%s": %s, ignoring file...`+"\n", filePath, err)
			continue
		}

		if cacheEntry.IsExpired() {
			numDeletedFiles += deleteCacheEntry(filePath)
		}
	}

	emptyDirectories := collectEmptyDirectoriesInPath(path)

	for _, directoryPath := range emptyDirectories {
		log.Printf(`deleting empty directory "%s"...`+"\n", directoryPath)

		if !flagDryRun {
			err := os.Remove(directoryPath)
			if err != nil {
				log.Printf(`error with directory "%s": %s`+"\n", directoryPath, err)
			}
		}
	}

	return numDeletedFiles
}

func collectCacheEntryFilesInPath(path string) []string {
	var cacheEntryFiles []string

	err := filepath.Walk(path, func(path string, _ os.FileInfo, _ error) error {
		if filepath.Base(path) == cacheEntryFileName {
			cacheEntryFiles = append(cacheEntryFiles, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return cacheEntryFiles
}

func deleteCacheEntry(path string) int {
	files := findAdjacentFiles(path)

	numDeletedFiles := 0

	for _, file := range files {
		log.Printf(`deleting file "%s"...`+"\n", file)
		numDeletedFiles++

		if !flagDryRun {
			err := os.Remove(file)
			if err != nil {
				log.Printf(`could not delete file "%s", because: %s`+"\n", file, err)
			}
		}
	}

	return numDeletedFiles
}

func findAdjacentFiles(path string) []string {
	dir := filepath.Dir(path)
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Printf(`could not list files in directory "%s", because: %s`, dir, err)
		return nil
	}

	var files []string
	for _, fileInfo := range fileInfos {
		// ignore sub directories
		if fileInfo.IsDir() {
			continue
		}

		file := filepath.Join(dir, fileInfo.Name())
		files = append(files, file)
	}

	return files
}

func collectEmptyDirectoriesInPath(path string) []string {
	var directories []string

	err := filepath.Walk(path, func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() && isEmptyDirectory(path) {
			directories = append(directories, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return directories
}

func isEmptyDirectory(name string) bool {
	f, err := os.Open(filepath.Clean(name))
	if err != nil {
		return false
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	_, err = f.Readdirnames(1)
	return err == io.EOF
}
