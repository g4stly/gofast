package common

import (
	"github.com/g4stly/config"
	"flag"
	"fmt"
	"io"
	"os"
)

func init() {
	flag.Parse()
	Args = flag.Args()


	_, err := os.Stat(DotFileName)
	if err != nil {
		if err.(*os.PathError).Err == os.ErrNotExist {
			Fatal("init(): %v", err)
		}
		err = os.Mkdir(DotFileName, 0755)
	}
	if err != nil {
		Fatal("init(): %v", err)
	}

	configFileName := fmt.Sprintf("%v/config.json", DotFileName)
	Config, err = config.LoadFile(configFileName)
	if err != nil {
		Fatal("init(): %v", err)
	}
}

type Command interface {
	Help() int
	Exec([]string) int
}
// important stuffs
var Args []string
var Config map[string]interface{}
var verbose = flag.Bool("v", false, "verbose: print debug output")
var silent = flag.Bool("s", false, "silent: surpress all output")
var DotFileName = fmt.Sprintf("/home/%v/.gofast", os.Getenv("USER"))

/*
	logging stuffs
*/
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










