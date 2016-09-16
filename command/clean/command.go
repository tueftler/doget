package clean

import (
	"flag"
	"os"

	"github.com/tueftler/doget/command"
	"github.com/tueftler/doget/config"
	"github.com/tueftler/doget/dockerfile"
)

// CleanCommand allows to remove the vendor directory and all of its contents
type CleanCommand struct {
	command.Command
	flags *flag.FlagSet
}

// NewCommand creates new clean command instance
func NewCommand(name string) *CleanCommand {
	return &CleanCommand{flags: flag.NewFlagSet(name, flag.ExitOnError)}
}

// Run performs action of clean command
func (c *CleanCommand) Run(parser *dockerfile.Parser, args []string) error {
	target := config.Vendordir
	if _, err := os.Stat(target); nil == err {
		return os.RemoveAll(target)
	}

	return nil
}
