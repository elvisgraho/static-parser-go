package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/elvisgraho/static-parser-go/utils"
)

var (
	fetchUrl   string
	rootDir    string
	configPath string
)

func init() {
	flag.StringVar(&fetchUrl, "fetch", "", "Enable fetching of missing .js and .json files. Provide base URL https://example.com.")
	flag.StringVar(&rootDir, "d", "", "Directory path to read files from")
	flag.StringVar(&configPath, "config", "config.json", "Path to the configuration file")

	// Custom usage function
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usages:\n")
		fmt.Fprintln(flag.CommandLine.Output(), "This application parses all files in the specified directory (and subdirectories) to extract information based on regex patterns defined in the config.json file.")
		fmt.Fprintln(flag.CommandLine.Output(), "With the -fetch option, it also fetches additional files (.js and .json) that were not found locally from the specified domain.")
		fmt.Fprintln(flag.CommandLine.Output(), "\nOptions:")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if flag.NFlag() == 0 || rootDir == "" {
		flag.Usage()
		os.Exit(1)
	}

	parsedDir, parseDirErr := utils.EnsureParsedDir(rootDir)
	if parseDirErr != nil {
		fmt.Printf("Failed to create parsed directory: %v\n", parseDirErr)
		return
	}

	flag.Parse()

	// Load the configuration.
	config, configLoadErr := utils.LoadConfig(configPath)
	if configLoadErr != nil {
		log.Fatalf("Failed to load config: %v", configLoadErr)
	}

	filesMap, err := utils.PopulateFilesMap(rootDir, config.ExcludedExtensions)
	if err != nil {
		log.Fatalf("Failed to populate files map: %v", err)
	}

	if fetchUrl != "" {
		if !strings.HasSuffix(fetchUrl, "/") {
			fetchUrl += "/"
		}
		utils.FindTsToFetch(filesMap, fetchUrl, rootDir)
	} else {
		matchesMap := utils.ProcessFilesMap(filesMap, config.ParsingJobs)

		// De-duplicate matches for each key in matchesMap
		for key, matches := range matchesMap {
			matches = utils.RemoveDuplicates(matches)
			matches = utils.FilterUnwantedExtensions(matches, config.UnwantedOutExtensions)
			matches = utils.FilterUnwantedExtensions(matches, config.UnwantedStrings)
			matches = utils.AlphabeticallySort(matches)
			matchesMap[key] = matches
		}

		// Iterate over matchesMap and write findings to files
		for key, matches := range matchesMap {
			matches = utils.RemoveDuplicates(matches) // De-duplicate matches
			filePath := filepath.Join(parsedDir, key+".txt")
			err := utils.WriteFindingsToFile(filePath, matches)
			if err != nil {
				fmt.Printf("Failed to write findings for %s: %v\n", key, err)
			}
		}
	}

}
