package cli

import (
	"github.com/urfave/cli/v2"
	// log "github.com/sirupsen/logrus"
)

var migrateDBCommand = &cli.Command{
	Name:     "database",
	// Category: "migrate",
	Usage:    "migrate database",
	Action:   func(c *cli.Context) error {
		fin, err := getFin(c)
		if err != nil {
			return err
		}

		return fin.MigrateDB()
	},
}

var migrateFileInfoCommand = &cli.Command{
	Name:     "fileinfo",
	// Category: "migrate",
	Usage:    "updates fileinfo with new fields",
	Action:   func(c *cli.Context) error {
		fin, err := getFin(c)
		if err != nil {
			return err
		}

		return fin.RebuildFileInfo()
	},
}

var migrateCommand = &cli.Command{
	Name:    "update",
	Usage:   "updates db and fileinfo table between versions",
	Subcommands: []*cli.Command{
		migrateDBCommand,
		migrateFileInfoCommand,
	},
}
