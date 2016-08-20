package main

import (
	"fmt"
	"github.com/tueftler/doget/dockerfile"
	"path/filepath"
	"strings"
)

func comment(value string) {
	fmt.Printf("# %s\n", strings.Replace(value, "\n", "\n# ", -1))
}

func instruction(instruction, value string) {
	fmt.Printf("%s %s\n\n", instruction, strings.Replace(value, "\n", "\\\n", -1))
}

func write(file *dockerfile.Dockerfile, base string) error {
	for _, statement := range file.Statements {
		switch statement.(type) {
		case *Include:
			var path string

			path, err := fetch(statement.(*Include).Reference);
			if err != nil {
				return err
			}

			var included dockerfile.Dockerfile
			if err := parser.ParseFile(filepath.Join(path, "Dockerfile.in"), &included); err != nil {
				return err
			}

			if included.From.Image != file.From.Image {
				return fmt.Errorf(
					"Include %s inherits from %s, which is incompatible with %s",
					statement.(*Include).Reference,
					included.From.Image,
					file.From.Image,
				)
			}

			comment("Included from " + statement.(*Include).Reference)
			write(&included, path + "/")
			break

		// Retain comments
		case *dockerfile.Comment:
			comment(statement.(*dockerfile.Comment).Lines)
			break

		// Builtin Docker instructions
		case *dockerfile.Maintainer:
			instruction("MAINTAINER", statement.(*dockerfile.Maintainer).Name)
			break
		case *dockerfile.Run:
			instruction("RUN", statement.(*dockerfile.Run).Command)
			break
		case *dockerfile.Label:
			instruction("LABEL", statement.(*dockerfile.Label).Pairs)
			break
		case *dockerfile.Expose:
			instruction("EXPOSE", statement.(*dockerfile.Expose).Ports)
			break
		case *dockerfile.Env:
			instruction("ENV", statement.(*dockerfile.Env).Pairs)
			break
		case *dockerfile.Add:
			instruction("ADD", base + statement.(*dockerfile.Add).Paths)
			break
		case *dockerfile.Copy:
			instruction("COPY", base + statement.(*dockerfile.Copy).Paths)
			break
		case *dockerfile.Entrypoint:
			instruction("ENTRYPOINT", statement.(*dockerfile.Entrypoint).CmdLine)
			break
		case *dockerfile.Volume:
			instruction("VOLUME", statement.(*dockerfile.Volume).Names)
			break
		case *dockerfile.User:
			instruction("USER", statement.(*dockerfile.User).Name)
			break
		case *dockerfile.Workdir:
			instruction("WORKDIR", statement.(*dockerfile.Workdir).Path)
			break
		case *dockerfile.Arg:
			instruction("ARG", statement.(*dockerfile.Arg).Name)
			break
		case *dockerfile.Onbuild:
			instruction("ONBUILD", statement.(*dockerfile.Onbuild).Instruction)
			break
		case *dockerfile.Stopsignal:
			instruction("STOPSIGNAL", statement.(*dockerfile.Stopsignal).Signal)
			break
		case *dockerfile.Healthcheck:
			instruction("HEALTHCHECK", statement.(*dockerfile.Healthcheck).Command)
			break
		case *dockerfile.Shell:
			instruction("SHELL", statement.(*dockerfile.Shell).CmdLine)
			break
		case *dockerfile.Cmd:
			instruction("CMD", statement.(*dockerfile.Cmd).CmdLine)
			break
		}
	}
	return nil
}

func transform(file *dockerfile.Dockerfile) error {
	instruction("FROM", file.From.Image)
	return write(file, "")
}