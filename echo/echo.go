package echo

import (
	"github.com/g4stly/gofast/common"
	"strings"
)

type command struct {
}

func (self *command) Exec(args []string) int {
	common.Out("%v\n", strings.Join(args, " "))
	return 0
}

func (self *command) Help() int {
	common.Out("echo <arg1> <arg2> ... <argN>\n[ECHO]: echos back the arguments you give it\n")
	return 0
}

func New() common.Command {
	cmd := command{}
	return &cmd
}
