package cli

import (
	"github.com/urfave/cli/v2"
	// log "github.com/sirupsen/logrus"
)

var statsCommand = &cli.Command{
	Name:    "stats",
	Usage:   "stats on database",
	Action:  func(c *cli.Context) error {
		fin, err := getFin(c)
		if err != nil {
			return err
		}

		return fin.Stats()
	},
}
