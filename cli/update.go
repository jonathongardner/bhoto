package cli

import (
	"github.com/urfave/cli/v2"
	// log "github.com/sirupsen/logrus"
)

var rebuildCommand = &cli.Command{
	Name:    "rebuild-fileinfo",
	Usage:   "rebuild fileinfo table",
	Action:  func(c *cli.Context) error {
		fin, err := getFin(c)
		if err != nil {
			return err
		}

		return fin.RebuildFileInfo()
	},
}
