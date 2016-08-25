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

func Test_emitting_maintainer(t *testing.T) {
	assertEmitted("MAINTAINER x\n\n", &Maintainer{Line: 1, Name: "x"}, t)
}

func Test_emitting_run(t *testing.T) {
	assertEmitted("RUN apt-get update\n\n", &Run{Line: 1, Command: "apt-get update"}, t)
}

func Test_emitting_multiline_run(t *testing.T) {
	assertEmitted("RUN apt-get -y install\\\n  doget\n\n", &Run{Line: 1, Command: "apt-get -y install\n  doget"}, t)
}

func Test_emitting_label(t *testing.T) {
	assertEmitted("LABEL key=value\n\n", &Label{Line: 1, Pairs: "key=value"}, t)
}

func Test_emitting_expose(t *testing.T) {
	assertEmitted("EXPOSE 80\n\n", &Expose{Line: 1, Ports: "80"}, t)
}

func Test_emitting_env(t *testing.T) {
	assertEmitted("ENV HTTP_PROXY\n\n", &Env{Line: 1, Pairs: "HTTP_PROXY"}, t)
}

func Test_emitting_add(t *testing.T) {
	assertEmitted("ADD source target\n\n", &Add{Line: 1, Paths: "source target"}, t)
}

func Test_emitting_copy(t *testing.T) {
	assertEmitted("COPY source target\n\n", &Copy{Line: 1, Paths: "source target"}, t)
}

func Test_emitting_entrypoint(t *testing.T) {
	assertEmitted("ENTRYPOINT cmd\n\n", &Entrypoint{Line: 1, CmdLine: "cmd"}, t)
}

func Test_emitting_volume(t *testing.T) {
	assertEmitted("VOLUME path\n\n", &Volume{Line: 1, Names: "path"}, t)
}

func Test_emitting_user(t *testing.T) {
	assertEmitted("USER name\n\n", &User{Line: 1, Name: "name"}, t)
}

func Test_emitting_workdir(t *testing.T) {
	assertEmitted("WORKDIR path\n\n", &Workdir{Line: 1, Path: "path"}, t)
}

func Test_emitting_arg(t *testing.T) {
	assertEmitted("ARG name\n\n", &Arg{Line: 1, Name: "name"}, t)
}

func Test_emitting_onbuild(t *testing.T) {
	assertEmitted("ONBUILD ADD source target\n\n", &Onbuild{Line: 1, Instruction: "ADD source target"}, t)
}

func Test_emitting_stopsignal(t *testing.T) {
	assertEmitted("STOPSIGNAL 11\n\n", &Stopsignal{Line: 1, Signal: "11"}, t)
}

func Test_emitting_healthcheck(t *testing.T) {
	assertEmitted("HEALTHCHECK cmd\n\n", &Healthcheck{Line: 1, Command: "cmd"}, t)
}

func Test_emitting_shell(t *testing.T) {
	assertEmitted("SHELL bash\n\n", &Shell{Line: 1, CmdLine: "bash"}, t)
}

func Test_emitting_cmd(t *testing.T) {
	assertEmitted("CMD /bin/bash\n\n", &Cmd{Line: 1, CmdLine: "/bin/bash"}, t)
}
