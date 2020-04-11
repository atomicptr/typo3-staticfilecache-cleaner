package cli

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	flagDryRun = false
)

var rootCommand = &cobra.Command{
	Use: "typo3-staticfilecache-cleaner",
	Run: func(cmd *cobra.Command, paths []string) {
		// no paths specified
		if len(paths) == 0 {
			log.Println("No paths specified.")
			os.Exit(1)
		}

		// check if all paths exist and are directories
		for _, path := range paths {
			fileInfo, err := os.Stat(path)

			if os.IsNotExist(err) {
				log.Printf(`"%s" does not exist.`, path)
				os.Exit(1)
			}

			if !fileInfo.IsDir() {
				log.Printf(`"%s" is not a directory.`, path)
				os.Exit(1)
			}
		}

		cleanPaths(paths)
	},
}

func init() {
	rootCommand.PersistentFlags().BoolVar(&flagDryRun, "dry-run", false, "do not delete anything")
}

func Run() error {
	return rootCommand.Execute()
}
