package utils

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func FilterUnwantedExtensions(paths []string, unwantedExtensions []string) []string {
	var filtered []string
	for _, path := range paths {
		// Split the path at " | " and only consider the first part
		parts := strings.Split(path, " | ")
		pathBeforePipe := strings.ToLower(parts[0]) // Convert to lowercase for case-insensitive comparison

		exclude := false
		for _, ext := range unwantedExtensions {
			if strings.Contains(pathBeforePipe, ext) {
				exclude = true
				break // No need to check other extensions
			}
		}
		if !exclude {
			filtered = append(filtered, path)
		}
	}
	return filtered
}

func RemoveDuplicates(slice []string) []string {
	unique := make(map[string]bool)
	var result []string

	for _, item := range slice {
		if _, found := unique[item]; !found {
			unique[item] = true
			result = append(result, item)
		}
	}

	return result
}

func AlphabeticallySort(slice []string) []string {
	// Make a copy of the slice to avoid modifying the original slice
	sortedSlice := make([]string, len(slice))
	copy(sortedSlice, slice)

	// Sort the copy of the slice
	sort.Strings(sortedSlice)

	return sortedSlice
}

func EnsureParsedDir(dir string) (string, error) {
	// Clean the directory path
	dir = filepath.Clean(dir)

	// Get the parent directory and the base name of the input directory
	parentDir := filepath.Dir(dir)
	baseName := filepath.Base(dir)

	// Construct the parsed directory path by appending '_parsed' to the base name
	parsedDir := filepath.Join(parentDir, baseName+"_parsed")

	// Create the parsed directory
	if err := os.MkdirAll(parsedDir, os.ModePerm); err != nil {
		return "", err
	}

	return parsedDir, nil
}

// Function to write findings to a file
func WriteFindingsToFile(filePath string, findings []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, finding := range findings {
		if _, err := file.WriteString(finding + "\n"); err != nil {
			return err
		}
	}
	return nil
}
