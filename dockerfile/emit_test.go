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
