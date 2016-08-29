package use

import (
	"reflect"
	"strings"
	"testing"

	"github.com/tueftler/doget/config"
	"github.com/tueftler/doget/dockerfile"
)

func assertEqual(expect, actual interface{}, t *testing.T) {
	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("Items not equal:\nexpected %q\nhave     %q\n", expect, actual)
	}
}

func mustParse(input string) *Statement {
	var file dockerfile.Dockerfile

	parser := dockerfile.NewParser().Extend("USE", New(config.Default().Repositories).Extension)
	if err := parser.Parse(strings.NewReader(input), &file); err != nil {
		panic(err)
	}

	return file.Statements[0].(*Statement)
}

func Test_reference(t *testing.T) {
	assertEqual("github.com/thekid/trait", mustParse("USE github.com/thekid/trait").Reference, t)
}

func Test_reference_with_trailing_space(t *testing.T) {
	assertEqual("github.com/thekid/trait", mustParse("USE github.com/thekid/trait ").Reference, t)
}

func Test_reference_with_leading_space(t *testing.T) {
	assertEqual("github.com/thekid/trait", mustParse("USE  github.com/thekid/trait").Reference, t)
}

func Test_origin_host(t *testing.T) {
	origin, err := mustParse("USE github.com/thekid/trait").Origin()
	if err != nil {
		t.Error(err.Error())
		return
	}
	assertEqual("github.com", origin.Host, t)
}

func Test_origin_illegal_host(t *testing.T) {
	_, err := mustParse("USE example.com/thekid/trait").Origin()
	if err == nil {
		t.Error("Expected an error, have non")
		return
	}
	assertEqual("No repository example.com", err.Error(), t)
}

func Test_origin_vendor(t *testing.T) {
	origin, err := mustParse("USE github.com/thekid/trait").Origin()
	if err != nil {
		t.Error(err.Error())
		return
	}
	assertEqual("thekid", origin.Vendor, t)
}

func Test_origin_name(t *testing.T) {
	origin, err := mustParse("USE github.com/thekid/trait").Origin()
	if err != nil {
		t.Error(err.Error())
		return
	}
	assertEqual("trait", origin.Name, t)
}

func Test_origin_without_dir(t *testing.T) {
	origin, err := mustParse("USE github.com/thekid/trait").Origin()
	if err != nil {
		t.Error(err.Error())
		return
	}
	assertEqual("", origin.Dir, t)
}

func Test_origin_with_dir(t *testing.T) {
	origin, err := mustParse("USE github.com/thekid/trait/dir").Origin()
	if err != nil {
		t.Error(err.Error())
		return
	}
	assertEqual("dir", origin.Dir, t)
}

func Test_origin_with_subdir(t *testing.T) {
	origin, err := mustParse("USE github.com/thekid/trait/sub/dir").Origin()
	if err != nil {
		t.Error(err.Error())
		return
	}
	assertEqual("sub/dir", origin.Dir, t)
}

func Test_origin_without_version_defaults_to_master(t *testing.T) {
	origin, err := mustParse("USE github.com/thekid/trait").Origin()
	if err != nil {
		t.Error(err.Error())
		return
	}
	assertEqual("master", origin.Version, t)
}

func Test_origin_with_version(t *testing.T) {
	origin, err := mustParse("USE github.com/thekid/trait:v1.0.0").Origin()
	if err != nil {
		t.Error(err.Error())
		return
	}
	assertEqual("v1.0.0", origin.Version, t)
}
