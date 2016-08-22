package command

import (
	"github.com/tueftler/doget/config"
	"github.com/tueftler/doget/dockerfile"
)

type Command interface {
	Init(name string)
	Run(config *config.Configuration, parser *dockerfile.Parser, args []string) error
}
