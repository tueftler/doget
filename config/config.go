package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Source       string
	Repositories map[string]map[string]string `yaml:"repositories"`
}

var (
	search = []func() string{
		func() string { return ".doget.yml" },
		func() string { return filepath.Join(filepath.Dir(os.Args[0]), ".doget.yml") },
		func() string { return filepath.Join(os.Getenv("HOME"), ".doget.yml") },
		func() string { return filepath.Join(os.Getenv("APPDATA"), "Doget", "config.yml") },
	}
)

// Default configuration loaded from search path
func Default() (result *Configuration, err error) {
	return From(SearchPath()...)
}

// Returns search path
func SearchPath() []string {
	result := make([]string, len(search))
	for i, path := range search {
		result[i] = path()
	}

	return result
}

// Read configuration from given sources
func From(sources ...string) (result *Configuration, err error) {
	result = &Configuration{Source: "", Repositories: make(map[string]map[string]string)}

	parsed := make(map[string]bool)

	for _, file := range sources {
		_, err = os.Stat(file)
		if err != nil {
			continue
		}

		info, err := filepath.Abs(file)
		if err != nil {
			continue
		}

		if _, ok := parsed[info]; ok {
			continue
		}

		parsed[info] = true
		parsedFile, err := FromFile(file)
		if err != nil {
			return nil, err
		}

		// Merge
		result.Source += ";" + parsedFile.Source
		for host, config := range parsedFile.Repositories {
			result.Repositories[host] = config
		}
	}

	result.Source = strings.TrimLeft(result.Source, ";")
	return result, nil
}

// Read configuration from a given file
func FromFile(filename string) (result *Configuration, err error) {
	result = &Configuration{Source: "", Repositories: make(map[string]map[string]string)}

	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(source, &result)
	if err != nil {
		return nil, err
	}

	result.Source = filename
	return result, nil
}
