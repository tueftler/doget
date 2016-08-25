package dockerfile

import (
	"fmt"
	"io"
	"strings"
)

// EmitComment writes a comment
func EmitComment(out io.Writer, value string) {
	fmt.Fprintf(out, "# %s\n", strings.Replace(value, "\n", "\n# ", -1))
}

// EmitInstruction writes an instruction
func EmitInstruction(out io.Writer, instruction, value string) {
	fmt.Fprintf(out, "%s %s\n\n", instruction, strings.Replace(value, "\n", "\\\n", -1))
}

// Emit writes comments
func (c *Comment) Emit(out io.Writer) {
	EmitComment(out, c.Lines)
}

// Emit writes FROM instruction
func (f *From) Emit(out io.Writer) {
	EmitInstruction(out, "FROM", f.Image)
}

// Emit writes MAINTAINER instructions
func (m *Maintainer) Emit(out io.Writer) {
	EmitInstruction(out, "MAINTAINER", m.Name)
}

// Emit writes RUN instructions
func (r *Run) Emit(out io.Writer) {
	EmitInstruction(out, "RUN", r.Command)
}

// Emit writes LABEL instructions
func (l *Label) Emit(out io.Writer) {
	EmitInstruction(out, "LABEL", l.Pairs)
}

// Emit writes EXPOSE instructions
func (e *Expose) Emit(out io.Writer) {
	EmitInstruction(out, "EXPOSE", e.Ports)
}

// Emit writes ENV instructions
func (e *Env) Emit(out io.Writer) {
	EmitInstruction(out, "ENV", e.Pairs)
}

// Emit writes ADD instructions
func (a *Add) Emit(out io.Writer) {
	EmitInstruction(out, "ADD", a.Paths)
}

// Emit writes COPY instructions
func (c *Copy) Emit(out io.Writer) {
	EmitInstruction(out, "COPY", c.Paths)
}

// Emit writes ENTRYPOINT instructions
func (e *Entrypoint) Emit(out io.Writer) {
	EmitInstruction(out, "ENTRYPOINT", e.CmdLine)
}

// Emit writes VOLUME instructions
func (v *Volume) Emit(out io.Writer) {
	EmitInstruction(out, "VOLUME", v.Names)
}

// Emit writes USER instructions
func (u *User) Emit(out io.Writer) {
	EmitInstruction(out, "USER", u.Name)
}

// Emit writes WORKDIR instructions
func (w *Workdir) Emit(out io.Writer) {
	EmitInstruction(out, "WORKDIR", w.Path)
}

// Emit writes ARG instructions
func (a *Arg) Emit(out io.Writer) {
	EmitInstruction(out, "ARG", a.Name)
}

// Emit writes ONBUILD instructions
func (o *Onbuild) Emit(out io.Writer) {
	EmitInstruction(out, "ONBUILD", o.Instruction)
}

// Emit writes STOPSIGNAL instructions
func (s *Stopsignal) Emit(out io.Writer) {
	EmitInstruction(out, "STOPSIGNAL", s.Signal)
}

// Emit writes HEALTHCHECK instructions
func (h *Healthcheck) Emit(out io.Writer) {
	EmitInstruction(out, "HEALTHCHECK", h.Command)
}

// Emit writes SHELL instructions
func (s *Shell) Emit(out io.Writer) {
	EmitInstruction(out, "SHELL", s.CmdLine)
}

// Emit writes CMD instructions
func (c *Cmd) Emit(out io.Writer) {
	EmitInstruction(out, "CMD", c.CmdLine)
}
