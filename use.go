package main

import (
	"github.com/tueftler/doget/dockerfile"
)

type Use struct {
	Line      int
	Reference string
}

func use(file *dockerfile.Dockerfile, line int, tokens *dockerfile.Tokens) dockerfile.Statement {
	return &Use{Line: line, Reference: tokens.NextLine()}
}
