package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/elvisgraho/static-parser-go/models"
)

func PopulateFilesMap(root string, excludedExtensions []string) (map[string]string, error) {
	filesMap := make(map[string]string)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if !IsExcludedExtension(ext, excludedExtensions) {
				contentBytes, err := os.ReadFile(path)
				if err != nil {
					fmt.Printf("Error reading file %s: %v\n", path, err)
					return nil // Continue walking the directory tree even if this file fails
				}
				filesMap[path] = string(contentBytes)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return filesMap, nil
}

// LoadConfig reads and parses the configuration from a JSON file.
func LoadConfig(filePath string) (*models.Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config models.Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// IsExcludedExtension checks if a file extension is in the list of excluded extensions
func IsExcludedExtension(ext string, excludedExtensions []string) bool {
	for _, excludedExt := range excludedExtensions {
		if ext == excludedExt {
			return true
		}
	}
	return false
}
