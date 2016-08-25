package command

import (
	"github.com/tueftler/doget/dockerfile"
)

type Command interface {
	Init(name string)
	Run(parser *dockerfile.Parser, args []string) error
}
