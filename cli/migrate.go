package cli

import (
	"github.com/urfave/cli/v2"
	// log "github.com/sirupsen/logrus"
)

var migrateCommand = &cli.Command{
	Name:    "migrate",
	Usage:   "migrate database",
	Action:  func(c *cli.Context) error {
		fin, err := getFin(c)
		if err != nil {
			return err
		}

		return fin.MigrateDB()
	},
}
