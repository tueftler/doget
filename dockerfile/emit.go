package dockerfile

import (
	"fmt"
	"io"
	"strings"
)

// Comment writes a comment
func WriteComment(out io.Writer, value string) {
	fmt.Fprintf(out, "# %s\n", strings.Replace(value, "\n", "\n# ", -1))
}

// Instruction writes an instruction
func WriteInstruction(out io.Writer, instruction, value string) {
	fmt.Fprintf(out, "%s %s\n\n", instruction, strings.Replace(value, "\n", "\\\n", -1))
}

func (c *Comment) Emit(out io.Writer) {
	WriteComment(out, c.Lines)
}

func (f *From) Emit(out io.Writer) {
	WriteInstruction(out, "FROM", f.Image)
}

func (m *Maintainer) Emit(out io.Writer) {
	WriteInstruction(out, "MAINTAINER", m.Name)
}

func (r *Run) Emit(out io.Writer) {
	WriteInstruction(out, "RUN", r.Command)
}

func (l *Label) Emit(out io.Writer) {
	WriteInstruction(out, "LABEL", l.Pairs)
}

func (e *Expose) Emit(out io.Writer) {
	WriteInstruction(out, "EXPOSE", e.Ports)
}

func (e *Env) Emit(out io.Writer) {
	WriteInstruction(out, "ENV", e.Pairs)
}

func (a *Add) Emit(out io.Writer) {
	WriteInstruction(out, "ADD", a.Paths)
}

func (c *Copy) Emit(out io.Writer) {
	WriteInstruction(out, "COPY", c.Paths)
}

func (e *Entrypoint) Emit(out io.Writer) {
	WriteInstruction(out, "ENTRYPOINT", e.CmdLine)
}

func (v *Volume) Emit(out io.Writer) {
	WriteInstruction(out, "VOLUME", v.Names)
}

func (u *User) Emit(out io.Writer) {
	WriteInstruction(out, "USER", u.Name)
}

func (w *Workdir) Emit(out io.Writer) {
	WriteInstruction(out, "WORKDIR", w.Path)
}

func (a *Arg) Emit(out io.Writer) {
	WriteInstruction(out, "ARG", a.Name)
}

func (o *Onbuild) Emit(out io.Writer) {
	WriteInstruction(out, "ONBUILD", o.Instruction)
}

func (s *Stopsignal) Emit(out io.Writer) {
	WriteInstruction(out, "STOPSIGNAL", s.Signal)
}

func (h *Healthcheck) Emit(out io.Writer) {
	WriteInstruction(out, "HEALTHCHECK", h.Command)
}

func (s *Shell) Emit(out io.Writer) {
	WriteInstruction(out, "SHELL", s.CmdLine)
}

func (c *Cmd) Emit(out io.Writer) {
	WriteInstruction(out, "CMD", c.CmdLine)
}
