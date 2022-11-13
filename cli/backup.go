package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jonathongardner/bhoto/fileInventory"
	"github.com/jonathongardner/bhoto/hub"

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

		log.Infof("Starting %v...", folder)

		ctx, cancel := context.WithCancel(context.Background())

		dbPath, err := getDatabasePath(c)
		if err != nil {
			return err
		}

		fin, err := fileInventory.NewFin()
		if err != nil {
			return err
		}
		defer fin.Close()
		go fin.StartDB(dbPath)

		// listen for ctrl + c and gracefully shutdown
		go func() {
			c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)

			<-c
			log.Info("Gracefully Shuting down...")
			cancel()
		}()


		h := hub.NewHub(folder, fin)
		go h.Process(ctx, maxNumberOfFileProcessors)

		return h.Wait()
	},
}
