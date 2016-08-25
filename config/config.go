package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Configuration holds the map of all repository informations parsed from all config files
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

// Empty configuration
func Empty() *Configuration {
	return &Configuration{Source: "", Repositories: make(map[string]map[string]string)}
}

// Default configuration supports github.com and bitbucket.org
func Default() *Configuration {
	return &Configuration{Source: "<default>", Repositories: map[string]map[string]string{
		"github.com": map[string]string{
			"url": "https://github.com/{{.Vendor}}/{{.Name}}/archive/{{.Version}}.zip",
		},
		"bitbucket.org": map[string]string{
			"url": "https://bitbucket.org/{{.Vendor}}/{{.Name}}/get/{{.Version}}.zip",
		},
	}}
}

// SearchPath Returns search path
func SearchPath() []string {
	result := make([]string, len(search))
	for i, path := range search {
		result[i] = path()
	}

	return result
}

// Merge Reads configuration from given sources
func (c *Configuration) Merge(sources ...string) (*Configuration, error) {
  return c.merge(false, sources...)
}

// MustMerge Reads configuration from given sources and yields an error if none exist
func (c *Configuration) MustMerge(sources ...string) (*Configuration, error) {
  return c.merge(true, sources...)
}

func (c *Configuration) merge(must bool, sources ...string) (*Configuration, error) {
	parsed := make(map[string]bool)
	for _, file := range sources {
		if _, err := os.Stat(file); err != nil {
			continue
		}

		path, err := filepath.Abs(file)
		if err != nil {
			continue
		}

		if _, ok := parsed[path]; ok {
			continue
		}

		parsed[path] = true
		parsedFile, err := FromFile(path)
		if err != nil {
			return nil, err
		}

		// Merge
		c.Source += ";" + parsedFile.Source
		for host, config := range parsedFile.Repositories {
			c.Repositories[host] = config
		}
	}

	if 0 == len(parsed) && must {
		return nil, fmt.Errorf("None of the given config files exist: %q", sources)
	}

	c.Source = strings.TrimLeft(c.Source, ";")
	return c, nil
}

// FromFile Reads configuration from a given file
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
