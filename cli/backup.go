package cli

import (
	"context"
	"os"
	"os/signal"
	"fmt"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
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
			Usage:   "Folder to look for photos to backup (default: current directory)",
			EnvVars: []string{"BOTO_FOLDER"},
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

		log.Infof("Starting %v...", folder)

		ctx, cancel := context.WithCancel(context.Background())

		// listen for ctrl + c and gracefully shutdown
		go func() {
			c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)

			<-c
			log.Info("Gracefully Shuting down...")
			cancel()
		}()
		g, gCtx := errgroup.WithContext(ctx) // gCtx

		g.Go(func() error {
			for i := 1; i < 5; i++ {
				log.Infof("Hello %v", i)
			}
			return nil
		})
		g.Go(func() error {
			i := 0
			for {
				select {
				case <- gCtx.Done():
					return nil
				default:
					if i >= 10 {
						return fmt.Errorf("You took to long!")
					}
					time.Sleep(1 * time.Second)
					i++
				}
			}

		})

		// now wait to see if any errors are raised, if one is raised than it will call cancel and end
		if err := g.Wait(); err != nil {
			return err
		}

		log.Info("...Closing")

		return nil
	},
}
