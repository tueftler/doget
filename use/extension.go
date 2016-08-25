package use

import (
	"github.com/tueftler/doget/dockerfile"
	"io"
)

// Statement represents a single USE statement
type Statement struct {
	Context   *Context
	Line      int
	Reference string
}

// Context represents the use context
type Context struct {
	Repositories map[string]map[string]string
}

// New creates a USE instruction backed by the given repositories
func New(repositories map[string]map[string]string) *Context {
	return &Context{Repositories: repositories}
}

// Emit writes the USE statement
func (s *Statement) Emit(out io.Writer) {
	dockerfile.EmitInstruction(out, "USE", s.Reference)
}

// Extension func for parser
func (c *Context) Extension(file *dockerfile.Dockerfile, line int, tokens *dockerfile.Tokens) dockerfile.Statement {
	return &Statement{Context: c, Line: line, Reference: tokens.NextLine()}
}
