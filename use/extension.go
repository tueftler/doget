package use

import (
	"github.com/tueftler/doget/dockerfile"
	"io"
)

type Statement struct {
	Ext       *Use
	Line      int
	Reference string
}

type Use struct {
	Repositories map[string]map[string]string
}

// New creates a Use instruction backed by the given repositories
func New(repositories map[string]map[string]string) *Use {
	return &Use{Repositories: repositories}
}

func (s *Statement) Emit(out io.Writer) {
	// TODO
}

// Extension func for parser
func (u *Use) Extension(file *dockerfile.Dockerfile, line int, tokens *dockerfile.Tokens) dockerfile.Statement {
	return &Statement{Ext: u, Line: line, Reference: tokens.NextLine()}
}
