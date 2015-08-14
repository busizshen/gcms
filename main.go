package main

import (
	"os"
	"runtime"

	"github.com/codegangsta/cli"

	"github.com/robvdl/gcms/cmd"
	"github.com/robvdl/gcms/config"
)

func init() {
	// As of Go 1.5 this will be the default so we won't need to do this anymore
	// Before Go 1.5, this actually defaults to 1 CPU unless you do this.
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	app := cli.NewApp()
	app.Name = config.AppName
	app.Usage = "Content management system"
	app.Version = config.AppVersion

	// list of available commands
	app.Commands = []cli.Command{
		cmd.CmdWeb,
	}

	app.Run(os.Args)
}
