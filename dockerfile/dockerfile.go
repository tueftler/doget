package dockerfile

type Statement interface{}

type Dockerfile struct {
	Source     string
	Statements []Statement
	From       *From
}

type Comment struct {
	Line  int
	Lines string
}

type From struct {
	Line  int
	Image string
}

type Maintainer struct {
	Line int
	Name string
}

type Run struct {
	Line    int
	Command string
}

type Label struct {
	Line  int
	Pairs string
}

type Expose struct {
	Line  int
	Ports string
}

type Env struct {
	Line  int
	Pairs string
}

type Add struct {
	Line  int
	Paths string
}

type Copy struct {
	Line  int
	Paths string
}

type Entrypoint struct {
	Line    int
	CmdLine string
}

type Volume struct {
	Line  int
	Names string
}

type User struct {
	Line int
	Name string
}

type Workdir struct {
	Line int
	Path string
}

type Arg struct {
	Line int
	Name string
}

type Onbuild struct {
	Line        int
	Instruction string
}

type Stopsignal struct {
	Line   int
	Signal string
}

type Healthcheck struct {
	Line    int
	Command string
}

type Shell struct {
	Line    int
	CmdLine string
}

type Cmd struct {
	Line    int
	CmdLine string
}
