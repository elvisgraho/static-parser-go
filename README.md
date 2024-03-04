# static-parser-go

static-parser-go is a simple app that parses web application files from a folder.  
You could for example use "Save All Resources" browser extension to download all website files and then
parse the with this app.  
  
It is a simple regex matching defined in config.json, but it also can fetch missing .js or .json files found in
documents.

## Installation

To install static-parser-go, you need to have Go installed on your machine.

```sh
go install github.com/elvisgraho/static-parser-go@latest
```

## Usage

To use static-parser-go, provide the following parameters:

* **-d** Directory to parse files from recursivly.
* **-fetch** Enable fetching of missing .js and .json files. Provide base URL ```-fetch 'https://example.com'```.
* **-config** Provide your own config path. ```-config './config.json'```

### Examples

#### Parse the folder

```sh
static-parser-go -d .\test\
```

#### Fetch the files

```sh
static-parser-go -d .\test\ -fetch "https://example.com/"
```

### License

static-parser-go is open-sourced software licensed under the MIT license.
