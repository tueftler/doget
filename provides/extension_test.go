package provides

import (
	"reflect"
	"strings"
	"testing"

	"github.com/tueftler/doget/dockerfile"
)

func assertEqual(expect, actual interface{}, t *testing.T) {
	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("Items not equal:\nexpected %q\nhave     %q\n", expect, actual)
	}
}

func mustParse(input string) *Statement {
	var file dockerfile.Dockerfile

	parser := dockerfile.NewParser().Extend("PROVIDES", Extension)
	if err := parser.Parse(strings.NewReader(input), &file); err != nil {
		panic(err)
	}

	return file.Statements[0].(*Statement)
}

func Test_list(t *testing.T) {
	assertEqual("thekid/trait:1.9", mustParse("PROVIDES thekid/trait:1.9").List, t)
}

func Test_one_image(t *testing.T) {
	assertEqual([]string{"thekid/trait:1.9"}, mustParse("PROVIDES thekid/trait:1.9").Images(), t)
}

func Test_two_images(t *testing.T) {
	assertEqual(
		[]string{"thekid/trait:1.9", "thekid/trait:latest"},
		mustParse("PROVIDES thekid/trait:1.9 thekid/trait:latest").Images(),
		t,
	)
}

func Test_images_with_whitespace(t *testing.T) {
	assertEqual(
		[]string{"thekid/trait:1.9", "thekid/trait:latest"},
		mustParse("PROVIDES thekid/trait:1.9   thekid/trait:latest").Images(),
		t,
	)
}

func Test_images_with_leading_and_trailing_whitespace(t *testing.T) {
	assertEqual(
		[]string{"thekid/trait:1.9", "thekid/trait:latest"},
		mustParse("PROVIDES  thekid/trait:1.9  thekid/trait:latest ").Images(),
		t,
	)
}
