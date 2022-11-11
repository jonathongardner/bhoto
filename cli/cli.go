package cli

import (
	"fmt"
	"os"
	"github.com/jonathongardner/wegyb/app"

	"github.com/urfave/cli/v2"
)


func Run() (error) {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println(c.App.Version)
	}

	app := &cli.App{
		Name: "bhoto",
		Version: app.Version,
		Usage: "We got your back!",
		Commands: []*cli.Command{
  		backupCommand,
  	},
	}
	return app.Run(os.Args)
}
