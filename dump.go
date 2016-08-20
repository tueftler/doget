package main

import (
	"fmt"
	"github.com/tueftler/doget/dockerfile"
)

func dump(file *dockerfile.Dockerfile) error {
	fmt.Println(file.Source, "{")
	for _, statement := range file.Statements {
		fmt.Printf("  %T %+v\n", statement, statement)
	}
	fmt.Println("}")
	return nil
}
