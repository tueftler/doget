package build

import (
	"errors"
	"github.com/tueftler/doget/dockerfile"
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

func Test_splitArgsRecognizesTransformArgWithoutValue(t *testing.T) {
	transformArgs, _ := split([]string{"-t", "foo:bar", "--doget-clean", "--no-cache=true"})
	assertEqual([]string{"-clean"}, transformArgs, t)
}

func Test_splitArgsRecognizesDockerBuildArgs(t *testing.T) {
	_, dockerArgs := split([]string{"-t", "foo:bar", "--no-cache=true", "--doget-no-cache=true", "."})
	assertEqual([]string{"-t", "foo:bar", "--no-cache=true", "."}, dockerArgs, t)
}

func Test_transformArgsAreEmptyWhenOnlyDockerArgsPresent(t *testing.T) {
	transformArgs, _ := split([]string{"--no-cache=true"})
	assertEqual([]string{}, transformArgs, t)
}

func Test_dockerArgsAreEmptyWhenOnlyTransformArgsPresent(t *testing.T) {
	_, dockerArgs := split([]string{"--doget-no-cache=true"})
	assertEqual([]string{}, dockerArgs, t)
}

type mock struct {
	executed bool
	err      error
}

func (m *mock) Build(args []string) error {
	m.executed = true
	if nil != m.err {
		return m.err
	}

	return nil
}

func (m *mock) Help() ([]byte, error) {
	return []byte{}, nil
}

func (m *mock) Init(name string) {
	// intentionally empty
}

func (m *mock) Run(parser *dockerfile.Parser, args []string) error {
	m.executed = true
	return m.err
}

func Test_dockerBuildNotExecutedWhenTransformFails(t *testing.T) {
	transform := &mock{executed: false, err: errors.New("an error")}
	docker := &mock{executed: false}
	clean := &mock{executed: false}
	NewCommand("build", transform, clean, docker).Run(dockerfile.NewParser(), []string{"."})
	assertEqual(false, docker.executed, t)
	assertEqual(false, clean.executed, t)
}

func Test_returnsTransformError(t *testing.T) {
	err := errors.New("an error")
	transform := &mock{executed: false, err: err}
	docker := &mock{executed: false}
	clean := &mock{executed: false}

	assertEqual(err, NewCommand("build", transform, clean, docker).Run(dockerfile.NewParser(), []string{"."}), t)
}

func Test_cleanNotExecutedWhenDockerBuildFails(t *testing.T) {
	transform := &mock{executed: false}
	docker := &mock{executed: false, err: errors.New("an error")}
	clean := &mock{executed: false}
	NewCommand("build", transform, clean, docker).Run(dockerfile.NewParser(), []string{"."})
	assertEqual(false, clean.executed, t)
}

func Test_returnsDockerBuildError(t *testing.T) {
	err := errors.New("an error")
	transform := &mock{executed: false}
	docker := &mock{executed: false, err: err}
	clean := &mock{executed: false}

	assertEqual(err, NewCommand("build", transform, clean, docker).Run(dockerfile.NewParser(), []string{"."}), t)
}

func Test_returnsCleanError(t *testing.T) {
	err := errors.New("an error")
	transform := &mock{executed: false}
	docker := &mock{executed: false}
	clean := &mock{executed: false, err: err}

	assertEqual(err, NewCommand("build", transform, clean, docker).Run(dockerfile.NewParser(), []string{"."}), t)
}

func Test_executedAllWhenNoneFails(t *testing.T) {
	transform := &mock{executed: false}
	docker := &mock{executed: false}
	clean := &mock{executed: false}
	NewCommand("build", transform, clean, docker).Run(dockerfile.NewParser(), []string{"."})
	assertEqual(true, transform.executed, t)
	assertEqual(true, docker.executed, t)
	assertEqual(true, clean.executed, t)
}

var showsUsage = []struct {
	args []string
}{
	{[]string{}},
	{[]string{"-help"}},
	{[]string{"--help"}},
}

func Test_showsUsage(t *testing.T) {
	for _, tt := range showsUsage {
		transform := &mock{executed: false}
		docker := &mock{executed: false}
		clean := &mock{executed: false}

		NewCommand("build", transform, clean, docker).Run(dockerfile.NewParser(), tt.args)
		assertEqual(false, transform.executed, t)
		assertEqual(false, docker.executed, t)
		assertEqual(false, clean.executed, t)
	}
}
