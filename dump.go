package main

import (
	"fmt"
	"github.com/tueftler/doget/dockerfile"
	"io"
)

func dump(out io.Writer, file *dockerfile.Dockerfile) error {
	fmt.Fprintln(out, file.Source, "{")
	for _, statement := range file.Statements {
		fmt.Fprintf(out, "  %T %+v\n", statement, statement)
	}
	fmt.Fprintln(out, "}")
	return nil
}
