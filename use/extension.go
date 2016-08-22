package use

import (
	"github.com/tueftler/doget/dockerfile"
)

type Statement struct {
	Line      int
	Reference string
}

// Extension func for parser
func Extension(file *dockerfile.Dockerfile, line int, tokens *dockerfile.Tokens) dockerfile.Statement {
	return &Statement{Line: line, Reference: tokens.NextLine()}
}
