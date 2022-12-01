package cli

import (
	"fmt"
	"os"

	"github.com/jonathongardner/bhoto/dirEntry"
	"github.com/jonathongardner/bhoto/routines"

	"github.com/urfave/cli/v2"
	log "github.com/sirupsen/logrus"
)

var backupCommand =  &cli.Command{
	Name:    "backup",
	Aliases: []string{"b"},
	Usage:   "backup photos in folder",
	Flags: []cli.Flag {
		&cli.StringFlag{
			Name:    "folder",
			Aliases: []string{"f"},
			Value:   "",
			DefaultText: "./",
			Usage:   "Folder to look for photos to backup",
			EnvVars: []string{"BOTO_FOLDER"},
		},
		&cli.IntFlag{
			Name:    "max",
			Aliases: []string{"m"},
			Value:   10,
			Usage:   "Max number of files to process",
			EnvVars: []string{"BOTO_MAX_FILES"},
			Action: func(ctx *cli.Context, v int) error {
				if 0 >= v {
					return fmt.Errorf("Flag max number of files to process %v must be greater than 0", v)
				}
				return nil
			},
		},
	},
	Action:  func(c *cli.Context) error {
		folder := c.String("folder")
		if folder == "" {
			var err error
			folder, err = os.Getwd()
			if err != nil {
				return err
			}
		}

		maxNumberOfFileProcessors := c.Int("max")

		routineController := routines.NewController(maxNumberOfFileProcessors)

		log.Infof("Starting %v...", folder)

		fin, err := getFin(c)
		if err != nil {
			return err
		}
		de, err := dirEntry.NewDirEntry(folder, fin)
		if err != nil {
			return err
		}

		routineController.GoBackground(fin)
		routineController.Go(de)

		return routineController.Wait()
	},
}
