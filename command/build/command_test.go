package build

import (
	"errors"
	"testing"

	"github.com/tueftler/doget/dockerfile"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type TestSuite struct{}

var _ = Suite(&TestSuite{})

func (s *TestSuite) Test_splitArgsRecognizesTransformArgs(c *C) {
	transformArgs, _ := split([]string{"--doget-no-cache=true", "-t", "foo:bar", "--no-cache=true"})
	c.Assert([]string{"-no-cache", "true"}, DeepEquals, transformArgs)
}

func (s *TestSuite) Test_splitArgsRecognizesTransformArgWithoutValue(c *C) {
	transformArgs, _ := split([]string{"-t", "foo:bar", "--doget-clean", "--no-cache=true"})
	c.Assert([]string{"-clean"}, DeepEquals, transformArgs)
}

func (s *TestSuite) Test_splitArgsRecognizesDockerBuildArgs(c *C) {
	_, dockerArgs := split([]string{"-t", "foo:bar", "--no-cache=true", "--doget-no-cache=true", "."})
	c.Assert([]string{"-t", "foo:bar", "--no-cache=true", "."}, DeepEquals, dockerArgs)
}

func (s *TestSuite) Test_transformArgsAreEmptyWhenOnlyDockerArgsPresent(c *C) {
	transformArgs, _ := split([]string{"--no-cache=true"})
	c.Assert([]string{}, DeepEquals, transformArgs)
}

func (s *TestSuite) Test_dockerArgsAreEmptyWhenOnlyTransformArgsPresent(c *C) {
	_, dockerArgs := split([]string{"--doget-no-cache=true"})
	c.Assert([]string{}, DeepEquals, dockerArgs)
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

func (s *TestSuite) Test_dockerBuildNotExecutedWhenTransformFails(c *C) {
	transform := &mock{executed: false, err: errors.New("an error")}
	docker := &mock{executed: false}
	clean := &mock{executed: false}
	NewCommand("build", transform, clean, docker).Run(dockerfile.NewParser(), []string{"."})
	c.Assert(false, Equals, docker.executed)
	c.Assert(false, Equals, clean.executed)
}

func (s *TestSuite) Test_returnsTransformError(c *C) {
	err := errors.New("an error")
	transform := &mock{executed: false, err: err}
	docker := &mock{executed: false}
	clean := &mock{executed: false}

	c.Assert(err, DeepEquals, NewCommand("build", transform, clean, docker).Run(dockerfile.NewParser(), []string{"."}))
}

func (s *TestSuite) Test_cleanNotExecutedWhenDockerBuildFails(c *C) {
	transform := &mock{executed: false}
	docker := &mock{executed: false, err: errors.New("an error")}
	clean := &mock{executed: false}
	NewCommand("build", transform, clean, docker).Run(dockerfile.NewParser(), []string{"."})
	c.Assert(false, Equals, clean.executed)
}

func (s *TestSuite) Test_returnsDockerBuildError(c *C) {
	err := errors.New("an error")
	transform := &mock{executed: false}
	docker := &mock{executed: false, err: err}
	clean := &mock{executed: false}

	c.Assert(err, DeepEquals, NewCommand("build", transform, clean, docker).Run(dockerfile.NewParser(), []string{"."}))
}

func (s *TestSuite) Test_returnsCleanError(c *C) {
	err := errors.New("an error")
	transform := &mock{executed: false}
	docker := &mock{executed: false}
	clean := &mock{executed: false, err: err}

	c.Assert(err, DeepEquals, NewCommand("build", transform, clean, docker).Run(dockerfile.NewParser(), []string{"."}))
}

func (s *TestSuite) Test_executedAllWhenNoneFails(c *C) {
	transform := &mock{executed: false}
	docker := &mock{executed: false}
	clean := &mock{executed: false}
	NewCommand("build", transform, clean, docker).Run(dockerfile.NewParser(), []string{"."})
	c.Assert(true, Equals, transform.executed)
	c.Assert(true, Equals, docker.executed)
	c.Assert(true, Equals, clean.executed)
}

var showsUsage = []struct {
	args []string
}{
	{[]string{}},
	{[]string{"-help"}},
	{[]string{"--help"}},
}

func (s *TestSuite) Test_showsUsage(c *C) {
	for _, tt := range showsUsage {
		transform := &mock{executed: false}
		docker := &mock{executed: false}
		clean := &mock{executed: false}

		NewCommand("build", transform, clean, docker).Run(dockerfile.NewParser(), tt.args)
		c.Check(false, Equals, transform.executed, Commentf("For %v", tt.args))
		c.Check(false, Equals, docker.executed, Commentf("For %v", tt.args))
		c.Check(false, Equals, clean.executed, Commentf("For %v", tt.args))
	}
}
