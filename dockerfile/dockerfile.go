package dockerfile

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Statement interface{}

type Dockerfile struct {
	Source     string
	Statements []Statement
	From       *From
}

type Comment struct {
	Statement
	Lines string
}

type From struct {
	Statement
	Image string
}

type Maintainer struct {
	Statement
	Name string
}

type Run struct {
	Statement
	Command string
}

type Label struct {
	Statement
	Pairs string
}

type Expose struct {
	Statement
	Ports string
}

type Env struct {
	Statement
	Pairs string
}

type Add struct {
	Statement
	Paths string
}

type Copy struct {
	Statement
	Paths string
}

type Entrypoint struct {
	Statement
	CmdLine string
}

type Volume struct {
	Statement
	Names string
}

type User struct {
	Statement
	Name string
}

type Workdir struct {
	Statement
	Path string
}

type Arg struct {
	Statement
	Name string
}

type Onbuild struct {
	Statement
	Instruction string
}

type Stopsignal struct {
	Statement
	Signal string
}

type Healthcheck struct {
	Statement
	Command string
}

type Shell struct {
	Statement
	CmdLine string
}

type Cmd struct {
	Statement
	CmdLine string
}

var (
	statements = map[string]func(file *Dockerfile, tokens *Tokens) Statement{
		"FROM": func(file *Dockerfile, tokens *Tokens) Statement {
			file.From = &From{Image: tokens.NextLine()}
			return file.From
		},
		"MAINTAINER": func(file *Dockerfile, tokens *Tokens) Statement {
			return &Maintainer{Name: tokens.NextLine()}
		},
		"RUN": func(file *Dockerfile, tokens *Tokens) Statement {
			return &Run{Command: tokens.NextLine()}
		},
		"CMD": func(file *Dockerfile, tokens *Tokens) Statement {
			return &Cmd{CmdLine: tokens.NextLine()}
		},
		"LABEL": func(file *Dockerfile, tokens *Tokens) Statement {
			return &Label{Pairs: tokens.NextLine()}
		},
		"EXPOSE": func(file *Dockerfile, tokens *Tokens) Statement {
			return &Expose{Ports: tokens.NextLine()}
		},
		"ENV": func(file *Dockerfile, tokens *Tokens) Statement {
			return &Env{Pairs: tokens.NextLine()}
		},
		"ADD": func(file *Dockerfile, tokens *Tokens) Statement {
			return &Add{Paths: tokens.NextLine()}
		},
		"COPY": func(file *Dockerfile, tokens *Tokens) Statement {
			return &Copy{Paths: tokens.NextLine()}
		},
		"ENTRYPOINT": func(file *Dockerfile, tokens *Tokens) Statement {
			return &Entrypoint{CmdLine: tokens.NextLine()}
		},
		"VOLUME": func(file *Dockerfile, tokens *Tokens) Statement {
			return &Volume{Names: tokens.NextLine()}
		},
		"USER": func(file *Dockerfile, tokens *Tokens) Statement {
			return &User{Name: tokens.NextLine()}
		},
		"WORKDIR": func(file *Dockerfile, tokens *Tokens) Statement {
			return &Workdir{Path: tokens.NextLine()}
		},
		"ARG": func(file *Dockerfile, tokens *Tokens) Statement {
			return &Arg{Name: tokens.NextLine()}
		},
		"ONBUILD": func(file *Dockerfile, tokens *Tokens) Statement {
			return &Onbuild{Instruction: tokens.NextLine()}
		},
		"STOPSIGNAL": func(file *Dockerfile, tokens *Tokens) Statement {
			return &Stopsignal{Signal: tokens.NextLine()}
		},
		"HEALTHCHECK": func(file *Dockerfile, tokens *Tokens) Statement {
			return &Healthcheck{Command: tokens.NextLine()}
		},
		"SHELL": func(file *Dockerfile, tokens *Tokens) Statement {
			return &Shell{CmdLine: tokens.NextLine()}
		},
		"#": func(file *Dockerfile, tokens *Tokens) Statement {
			return &Comment{Lines: tokens.NextComment()}
		},
	}
)

// Parses a dockerfile from a reader. Returns an error if
// an unknown token is encountered.
//
// See https://docs.docker.com/engine/reference/builder/
func Parse(input io.Reader, file *Dockerfile, source ...string) (err error) {
	tokens := NewTokens(input)

	for tokens.HasNext {
		token := tokens.NextToken()

		if "" == token {
			continue
		} else if statement, ok := statements[token]; ok {
			file.Statements = append(file.Statements, statement(file, tokens))
		} else {
			return fmt.Errorf("Cannot handle token `%s` on line %d", token, tokens.Line)
		}
	}

	if len(source) > 0 {
		file.Source = source[0]
	} else {
		file.Source = fmt.Sprintf("%T", input)
	}

	return nil
}

// Parses a dockerfile from a file. Returns an error if
// the file cannot be opened, is a directory or when parsing
// encounters an error
func ParseFile(name string, file *Dockerfile) (err error) {
	stat, err := os.Stat(name)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		return fmt.Errorf("The given file `%s` is a directory\n", name)
	}

	input, err := os.Open(name)
	if err != nil {
		return err
	}

	defer input.Close()
	return Parse(bufio.NewReader(input), file, name)
}
