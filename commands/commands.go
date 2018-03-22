package commands

import (
	"github.com/g4stly/gofast/common"
	"github.com/g4stly/gofast/create"
	"github.com/g4stly/gofast/echo"
)

func init() {
	Lookup = make(map[string]common.Command)
	Lookup["echo"] = echo.New()
	Lookup["new"] = create.New()
}

var Lookup map[string]common.Command
