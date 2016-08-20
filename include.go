package main

import (
	"github.com/tueftler/doget/dockerfile"
)

type Include struct {
	Line      int
	Reference string
}

func include(file *dockerfile.Dockerfile, line int, tokens *dockerfile.Tokens) dockerfile.Statement {
  return &Include{Line: line, Reference: tokens.NextLine()}
}