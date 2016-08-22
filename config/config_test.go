package config

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func assertEqual(expect, actual interface{}, t *testing.T) {
	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("Items not equal:\nexpected %q\nhave     %q\n", expect, actual)
	}
}

func configFile(content string) (*os.File, error) {
	file, err := ioutil.TempFile("", ".doget.yml")
	if err != nil {
		return nil, err
	}

	if _, err := file.Write([]byte(content)); err != nil {
		os.Remove(file.Name())
		return nil, err
	}
	if err := file.Close(); err != nil {
		os.Remove(file.Name())
		return nil, err
	}

	return file, nil
}

func Test_search_path(t *testing.T) {
	path := SearchPath()
	assertEqual(3, len(path), t)
}

func Test_can_parse_empty_file(t *testing.T) {
	file, err := configFile("")
	if err != nil {
		t.Errorf("Cannot create config file: %s", err.Error())
		return
	}
	defer os.Remove(file.Name())

	From(file.Name())
}

func Test_single_file_source(t *testing.T) {
	file, err := configFile("repositories:")
	if err != nil {
		t.Errorf("Cannot create config file: %s", err.Error())
		return
	}
	defer os.Remove(file.Name())

	config, _ := From(file.Name())
	assertEqual(file.Name(), config.Source, t)
}

func Test_multiple_file_sources(t *testing.T) {
	global, err := configFile("repositories:")
	if err != nil {
		t.Errorf("Cannot create config file: %s", err.Error())
		return
	}
	defer os.Remove(global.Name())

	user, err := configFile("repositories:")
	if err != nil {
		t.Errorf("Cannot create config file: %s", err.Error())
		return
	}
	defer os.Remove(user.Name())

	config, _ := From(global.Name(), user.Name())
	assertEqual(global.Name()+";"+user.Name(), config.Source, t)
}

func Test_single_file_repositories(t *testing.T) {
	file, err := configFile(`
repositories:
  github.com:
    url: https://github.com/...
`)
	if err != nil {
		t.Errorf("Cannot create config file: %s", err.Error())
		return
	}
	defer os.Remove(file.Name())

	config, err := From(file.Name())
	if err != nil {
		t.Errorf("Cannot parse config file: %s", err.Error())
		return
	}

	assertEqual("https://github.com/...", config.Repositories["github.com"]["url"], t)
}

func Test_adding_a_repository(t *testing.T) {
	global, err := configFile(`
repositories:
  github.com:
    url: https://github.com/...
`)
	if err != nil {
		t.Errorf("Cannot create config file: %s", err.Error())
		return
	}
	defer os.Remove(global.Name())

	user, err := configFile(`
repositories:
  example.com:
    url: https://example.com/...
`)
	if err != nil {
		t.Errorf("Cannot create config file: %s", err.Error())
		return
	}
	defer os.Remove(user.Name())

	config, _ := From(global.Name(), user.Name())
	assertEqual("https://github.com/...", config.Repositories["github.com"]["url"], t)
	assertEqual("https://example.com/...", config.Repositories["example.com"]["url"], t)
}

func Test_overwriting_a_repository(t *testing.T) {
	global, err := configFile(`
repositories:
  github.com:
    url: https://github.com/...
`)
	if err != nil {
		t.Errorf("Cannot create config file: %s", err.Error())
		return
	}
	defer os.Remove(global.Name())

	user, err := configFile(`
repositories:
  github.com:
    url: https://github.example.com/...
`)
	if err != nil {
		t.Errorf("Cannot create config file: %s", err.Error())
		return
	}
	defer os.Remove(user.Name())

	config, _ := From(global.Name(), user.Name())
	assertEqual("https://github.example.com/...", config.Repositories["github.com"]["url"], t)
}
