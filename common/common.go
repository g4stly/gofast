package common

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func init() {
	flag.Parse()
	Args = flag.Args()
}

type Command interface {
	Help() int
	Exec([]string) int
}

var Args []string
var verbose = flag.Bool("v", false, "verbose: print debug output")
var silent = flag.Bool("s", false, "silent: surpress all output")

func msg(w io.Writer, badge string, fmtstring string, args ...interface{}) {
	if *silent {
		return
	}
	if len(args) < 1 {
		w.Write([]byte(fmt.Sprintf("%v: %v\n", badge, fmtstring)))
	} else {
		w.Write([]byte(fmt.Sprintf("%v: %v\n", badge, fmt.Sprintf(fmtstring, args...))))
	}
}

func Log(fmtstring string, args ...interface{}) {
	if !*verbose {
		return
	}
	msg(os.Stdout, "DEBUG", fmtstring, args...)
}

func Out(fmtstring string, args ...interface{}) {
	msg(os.Stdout, "info", fmtstring, args...)
}

func Usage(fmtstring string, args ...interface{}) {
	msg(os.Stdout, "usage", fmtstring, args...)
}

func Fatal(fmtstring string, args ...interface{}) {
	msg(os.Stderr, "FATAL", fmtstring, args...)
	os.Exit(1)
}
