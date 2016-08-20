package main

import (
	"bytes"
	"fmt"
	"github.com/tueftler/doget/dockerfile"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func comment(out io.Writer, value string) {
	fmt.Fprintf(out, "# %s\n", strings.Replace(value, "\n", "\n# ", -1))
}

func instruction(out io.Writer, instruction, value string) {
	fmt.Fprintf(out, "%s %s\n\n", instruction, strings.Replace(value, "\n", "\\\n", -1))
}

func write(out io.Writer, file *dockerfile.Dockerfile, base string) error {
	for _, statement := range file.Statements {
		switch statement.(type) {
		case *Include:
			var path string
			reference := statement.(*Include).Reference

			path, err := fetch(reference, func(transferred, total int64) {
				percentage := float64(transferred) / float64(total)
				finished := int(percentage * float64(20))
				fmt.Fprintf(
					os.Stderr,
					"\r> Fetching %s: [%s%s] %.2fkB",
					reference,
					strings.Repeat("#", finished),
					strings.Repeat("_", 20-finished),
					float64(transferred)/float64(1024),
				)
			})
			fmt.Fprintf(os.Stderr, "\n")

			if err != nil {
				return err
			}

			var included dockerfile.Dockerfile
			if err := parse(filepath.Join(path, "Dockerfile.in"), &included); err != nil {
				return err
			}

			if included.From.Image != file.From.Image {
				return fmt.Errorf(
					"Include %s inherits from %s, which is incompatible with %s",
					reference,
					included.From.Image,
					file.From.Image,
				)
			}

			comment(out, "Included from "+reference)
			write(out, &included, path+"/")
			break

		// Retain comments
		case *dockerfile.Comment:
			comment(out, statement.(*dockerfile.Comment).Lines)
			break

		// Builtin Docker instructions
		case *dockerfile.Maintainer:
			instruction(out, "MAINTAINER", statement.(*dockerfile.Maintainer).Name)
			break
		case *dockerfile.Run:
			instruction(out, "RUN", statement.(*dockerfile.Run).Command)
			break
		case *dockerfile.Label:
			instruction(out, "LABEL", statement.(*dockerfile.Label).Pairs)
			break
		case *dockerfile.Expose:
			instruction(out, "EXPOSE", statement.(*dockerfile.Expose).Ports)
			break
		case *dockerfile.Env:
			instruction(out, "ENV", statement.(*dockerfile.Env).Pairs)
			break
		case *dockerfile.Add:
			instruction(out, "ADD", base+statement.(*dockerfile.Add).Paths)
			break
		case *dockerfile.Copy:
			instruction(out, "COPY", base+statement.(*dockerfile.Copy).Paths)
			break
		case *dockerfile.Entrypoint:
			instruction(out, "ENTRYPOINT", statement.(*dockerfile.Entrypoint).CmdLine)
			break
		case *dockerfile.Volume:
			instruction(out, "VOLUME", statement.(*dockerfile.Volume).Names)
			break
		case *dockerfile.User:
			instruction(out, "USER", statement.(*dockerfile.User).Name)
			break
		case *dockerfile.Workdir:
			instruction(out, "WORKDIR", statement.(*dockerfile.Workdir).Path)
			break
		case *dockerfile.Arg:
			instruction(out, "ARG", statement.(*dockerfile.Arg).Name)
			break
		case *dockerfile.Onbuild:
			instruction(out, "ONBUILD", statement.(*dockerfile.Onbuild).Instruction)
			break
		case *dockerfile.Stopsignal:
			instruction(out, "STOPSIGNAL", statement.(*dockerfile.Stopsignal).Signal)
			break
		case *dockerfile.Healthcheck:
			instruction(out, "HEALTHCHECK", statement.(*dockerfile.Healthcheck).Command)
			break
		case *dockerfile.Shell:
			instruction(out, "SHELL", statement.(*dockerfile.Shell).CmdLine)
			break
		case *dockerfile.Cmd:
			instruction(out, "CMD", statement.(*dockerfile.Cmd).CmdLine)
			break
		}
	}
	return nil
}

func transform(out io.Writer, file *dockerfile.Dockerfile) error {
	var transformed bytes.Buffer

	instruction(&transformed, "FROM", file.From.Image)
	if err := write(&transformed, file, ""); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Done\n\n")
	fmt.Fprintf(out, transformed.String())
	return nil
}
