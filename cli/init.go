package cli

import (
	"github.com/jonathongardner/bhoto/fileInventory"

	"github.com/urfave/cli/v2"
	// log "github.com/sirupsen/logrus"
)

var initCommand = &cli.Command{
	Name:    "init",
	Usage:   "setup database",
	Action:  func(c *cli.Context) error {
		path, err := getDatabasePath(c)
		if err != nil {
			return err
		}

		return fileInventory.SetupDB(path)
	},
}
