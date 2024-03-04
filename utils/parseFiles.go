package utils

import (
	"regexp"
	"runtime"
	"sync"

	"github.com/elvisgraho/static-parser-go/models"
)

func ProcessFilesMap(filesMap map[string]string, parsingJobs []models.ParsingJob) map[string][]string {
	matchesMap := make(map[string][]string)
	var mutex sync.Mutex

	for _, job := range parsingJobs {
		matchesMap[job.Key] = []string{} // Initialize the slice for each job key
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, runtime.NumCPU())

	for path, content := range filesMap {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire semaphore

		go func(path, content string) {
			defer wg.Done()
			for _, job := range parsingJobs {
				matches := applyRegex(content, job.RegexPattern)
				if len(matches) > 0 {
					mutex.Lock()
					matchesMap[job.Key] = append(matchesMap[job.Key], matches...)
					mutex.Unlock()
				}
			}
			<-semaphore // Release semaphore
		}(path, content)
	}

	wg.Wait()
	return matchesMap
}

func applyRegex(content string, pattern string) []string {
	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(content, -1) // Find all matches
	return matches
}
