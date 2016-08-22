package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Configuration struct {
	Source       string
	Repositories map[string]map[string]string `yaml:"repositories"`
}

var (
	search = []func() string{
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

	for _, file := range sources {
		_, err = os.Stat(file)
		if err != nil {
			continue
		}

		parsed, err := FromFile(file)
		if err != nil {
			return nil, err
		}

		// Merge
		result.Source += ";" + parsed.Source
		for host, config := range parsed.Repositories {
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
