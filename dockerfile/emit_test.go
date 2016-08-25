package dockerfile

import (
	"bytes"
	"reflect"
	"testing"
)

func assertEmitted(expect string, statement Statement, t *testing.T) {
	var buf bytes.Buffer

	statement.Emit(&buf)
	actual := buf.String()

	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("Items not equal:\nexpected %q\nhave     %q\n", expect, actual)
	}
}

func Test_emitting_from(t *testing.T) {
	assertEmitted("FROM debian:jessie\n\n", &From{Line: 1, Image: "debian:jessie"}, t)
}

func Test_emitting_comment(t *testing.T) {
	assertEmitted("# Test\n", &Comment{Line: 1, Lines: "Test"}, t)
}

func Test_emitting_multiline_comment(t *testing.T) {
	assertEmitted("# One\n# Two\n# Three\n", &Comment{Line: 1, Lines: "One\nTwo\nThree"}, t)
}

func Test_emitting_run(t *testing.T) {
	assertEmitted("RUN apt-get update\n\n", &Run{Line: 1, Command: "apt-get update"}, t)
}

func Test_emitting_multiline_run(t *testing.T) {
	assertEmitted("RUN apt-get -y install\\\n  doget\n\n", &Run{Line: 1, Command: "apt-get -y install\n  doget"}, t)
}
