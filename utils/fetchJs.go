package utils

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func FetchFile(path, baseURL string, rootDir string) {
	// Concatenate rootDir with "_fetched" to create the path for the fetched directory
	fetchedDirPath := filepath.Join(rootDir, "_fetched")

	var url string
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		url = path
	} else {
		// It's a relative path, so prepend the baseURL
		url = baseURL + path
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch %s: %v\n", url, err)
		return
	} else {
		log.Printf("Fetched %s\n", url)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to fetch %s: HTTP %d\n", url, resp.StatusCode)
		return
	}

	// Ensure the _fetched directory exists inside the rootDir
	errDir := os.MkdirAll(fetchedDirPath, 0755)
	if errDir != nil {
		log.Printf("Failed to create the _fetched directory: %v\n", errDir)
		return
	}

	// Extract the file name from the path and create the file in the _fetched directory
	filename := filepath.Base(url)
	filepath := filepath.Join(fetchedDirPath, filename)

	file, err := os.Create(filepath)
	if err != nil {
		log.Printf("Failed to create file %s: %v\n", filepath, err)
		return
	}
	defer file.Close()

	// Write the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Printf("Failed to write to file %s: %v\n", filepath, err)
		return
	}

	//log.Printf("Successfully fetched and saved %s\n", filepath)
}

func FindFilesToFetch(filesMap map[string]string, baseURL string, rootDir string) {
	tsRegex := regexp.MustCompile(`(?:https?:\/\/)?[\w/.\-]+\.js(?:on)?`)

	var foundPaths []string
	// Create a set of filenames from filesMap keys for quick lookup
	filenameSet := make(map[string]struct{})
	for key := range filesMap {
		filename := filepath.Base(key)
		filenameSet[filename] = struct{}{} // Add to set
	}

	for _, content := range filesMap {
		matches := tsRegex.FindAllString(content, -1)

		for _, match := range matches {
			matchFilename := filepath.Base(match)
			// Check if match's filename exists in the set of filenames from filesMap
			if _, exists := filenameSet[matchFilename]; !exists {
				// If the filename doesn't exist in the set, it's a new path that needs to be considered
				foundPaths = append(foundPaths, match)
			}
		}
	}

	foundPaths = RemoveDuplicates(foundPaths)

	var filteredPaths []string // This slice will hold the filtered paths

	for _, path := range foundPaths {
		// Check if the path starts with "https://" and does not contain the baseURL
		if strings.HasPrefix(path, "https://") && !strings.Contains(path, baseURL) {
			// If it's an https path without the baseURL, skip it
			continue
		}
		if strings.Contains(path, "node_modules") {
			// skip node modules
			continue
		}

		// If the path is either not https or is https with the baseURL, keep it
		filteredPaths = append(filteredPaths, path)
	}

	for _, path := range filteredPaths {
		FetchFile(path, baseURL, rootDir)
	}
}
