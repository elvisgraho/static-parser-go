package models

type Config struct {
	ParsingJobs           []ParsingJob `json:"parsingJobs"`
	ExcludedExtensions    []string     `json:"excludedExtensions"`
	UnwantedOutExtensions []string     `json:"unwantedOutExtensions"`
	UnwantedStrings       []string     `json:"unwantedStrings"`
}

type ParsingJob struct {
	Key          string `json:"key"`
	Title        string `json:"title"`
	RegexPattern string `json:"regexPattern"`
}
