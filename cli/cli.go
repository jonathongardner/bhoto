package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
	// "github.com/urfave/cli/v2/altsrc"
	log "github.com/sirupsen/logrus"
)
func getDatabasePath(c *cli.Context) (string, error) {
	folder := c.String("database")
	if folder == "" {
		var err error
		folder, err = os.Getwd()
		if err != nil {
			return "", err
		}
		return filepath.Join(folder, ".bhoto.sqlite"), nil
	}
	return folder, nil
}

func Run() (error) {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println(c.App.Version)
	}
	cli.VersionFlag = &cli.BoolFlag{
		Name: "version",
		Usage: "print the version",
	}

	flags := []cli.Flag {
		&cli.StringFlag{
			Name:    "database",
			Aliases: []string{"db"},
			DefaultText: "./.bhoto.sqlite",
			Usage:   "database file",
			EnvVars: []string{"BOTO_DB"},
		},
		// &cli.StringFlag{
		// 	Name:    "config",
		// 	Aliases: []string{"c"},
		// 	// DefaultText: "./.bhoto.yaml",
		// 	Usage:   "Load configuration from `FILE`",
		// 	EnvVars: []string{"BOTO_CONFIG"},
		// },
		&cli.BoolFlag{
			Name: "verbose",
			Aliases: []string{"v"},
			Usage: "logging level",
		},
	}


	app := &cli.App{
		Name: "bhoto",
		Version: "0.0.1",
		Usage: "We got your back!",
		Flags: flags,
		Before: func(c *cli.Context) error {
			// config := c.String("config")
			// if config != "" {
			// 	log.Infof("Loading config (%v)", config)
			// 	isc, err := altsrc.NewYamlSourceFromFile(config)
			// 	if err != nil {
			// 		return err
			// 	}
			// 	log.Infof("Test: %v", isc)
			//
			// 	err = altsrc.ApplyInputSourceValues(c, isc, initCommand.Flags)
			// 	if err != nil {
			// 		return err
			// 	}
			// }
			if c.Bool("verbose") {
				log.SetLevel(log.DebugLevel)
				log.Debug("Setting to debug...")
			}
			return nil
		},
		Commands: []*cli.Command{
			initCommand,
			backupCommand,
		},
	}
	return app.Run(os.Args)
}
