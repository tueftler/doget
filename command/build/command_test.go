package build

import (
	"reflect"
	"testing"
)

func assertEqual(expect, actual interface{}, t *testing.T) {
	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("Items not equal:\nexpected %q\nhave     %q\n", expect, actual)
	}
}

func Test_splitArgsRecognizesTransformArgs(t *testing.T) {
	transformArgs, _ := split([]string{"--doget-no-cache=true", "-t", "foo:bar", "--no-cache=true"})
	assertEqual([]string{"-no-cache", "true"}, transformArgs, t)
}

func Test_splitArgsRecognizesDockerBuildArgs(t *testing.T) {
	_, dockerArgs := split([]string{"--doget-no-cache=true", "-t", "foo:bar", "--no-cache=true"})
	assertEqual([]string{"build", "-t", "foo:bar", "--no-cache=true"}, dockerArgs, t)
}
