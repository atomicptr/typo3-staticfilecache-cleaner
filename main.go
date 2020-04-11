package main

import (
	"log"
	"os"

	"github.com/atomicptr/typo3-staticfilecache-cleaner/cli"
)

func main() {
	err := cli.Run()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
