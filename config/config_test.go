package config

import (
	"fmt"
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
	assertEqual(4, len(path), t)
}

func Test_empty_source(t *testing.T) {
	config := Empty()
	assertEqual("", config.Source, t)
}

func Test_default_source(t *testing.T) {
	config := Default()
	assertEqual("<default>", config.Source, t)
}

func Test_parse_only_nonexisting_files_does_return_error(t *testing.T) {
	_, err := Empty().Merge("doesNotExist", "doesNotExist2")
	assertEqual(fmt.Errorf("None of the given config files exist: [\"doesNotExist\" \"doesNotExist2\"]"), err, t)
}

func Test_parse_nonexisting_and_existing_file(t *testing.T) {
	file, err := configFile("repositories:")
	if err != nil {
		t.Errorf("Cannot create config file: %s", err.Error())
		return
	}
	defer os.Remove(file.Name())

	config, err := Empty().Merge("doesNotExist", file.Name())
	if nil != err {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	assertEqual(file.Name(), config.Source, t)
}

func Test_can_parse_empty_file(t *testing.T) {
	file, err := configFile("")
	if err != nil {
		t.Errorf("Cannot create config file: %s", err.Error())
		return
	}
	defer os.Remove(file.Name())

	Empty().Merge(file.Name())
}

func Test_single_file_source(t *testing.T) {
	file, err := configFile("repositories:")
	if err != nil {
		t.Errorf("Cannot create config file: %s", err.Error())
		return
	}
	defer os.Remove(file.Name())

	config, _ := Empty().Merge(file.Name())
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

	config, _ := Empty().Merge(global.Name(), user.Name())
	assertEqual(global.Name()+";"+user.Name(), config.Source, t)
}

func Test_same_file_sources_multiple_times(t *testing.T) {
	global, err := configFile("repositories:")
	if err != nil {
		t.Errorf("Cannot create config file: %s", err.Error())
		return
	}
	defer os.Remove(global.Name())

	config, _ := Empty().Merge(global.Name(), global.Name())
	assertEqual(global.Name(), config.Source, t)
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

	config, err := Empty().Merge(file.Name())
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

	config, _ := Empty().Merge(global.Name(), user.Name())
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

	config, _ := Empty().Merge(global.Name(), user.Name())
	assertEqual("https://github.example.com/...", config.Repositories["github.com"]["url"], t)
}
