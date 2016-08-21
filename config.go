package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Configuration struct {
	Source       string
	Repositories map[string]map[string]string `yaml:"repositories"`
}

func configuration(filename string) (result *Configuration, err error) {
	if "" == filename {
		filename = filepath.Join(filepath.Dir(os.Args[0]), ".doget.yml")
	}

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
