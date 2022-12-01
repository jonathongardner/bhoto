package routines

import (
  "fmt"
  "os"
  "os/signal"
  "strings"
  "syscall"

  log "github.com/sirupsen/logrus"
)

// struct to help control routines.
// will wait till all started routines finish to call `Wait`
// Might could have use sync error group but this will wait
// till all finish, rather than on first error
type Controller struct {
  //----running------
  runChan          chan Runner // true up false down
  runFinishedChan  chan bool // true up false down
  runCount         int // track number of open routines, that when 0 calls `cancel`
  //----waiting------
  backCountChan chan bool // true up false down
  backCount     int // track number of open routines, that we want to
  //----listeners------
  doneChan      chan struct{} // notify other of gracefully close
  finishedChan  chan struct{} // used for `Wait` is done
  //-----Pass errors-----
  errorChan     chan error
  errors        []string
}

type Runner interface {
	Run(rc *Controller) ([]Runner, error)
}
type BackgroundRunner interface {
	Run(rc *Controller) error
}


func NewController(maxNumberRunningThreads int) *Controller {
	// make dirEntryIn 1 buffer so we can go ahead and add the initialDir
	c := &Controller{
    runChan: make(chan Runner),
		runFinishedChan: make(chan bool),
		runCount: 0,
		backCountChan: make(chan bool),
		backCount: 0,
    doneChan: make(chan struct{}),
		finishedChan: make(chan struct{}),
    errorChan: make(chan error),
		errors: make([]string, 0),
	}

  go c.listenForCtrlC()
  go c.run(maxNumberRunningThreads)

	return c
}

func (c *Controller) listenForCtrlC() {
  ch := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
  signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

  <-ch
  log.Info("Gracefully Shuting down...")
  c.gracefullyClose()

  <-ch
  log.Info("Killing...")
  c.finish()
}

func (c *Controller) run(maxNumberRunningThreads int) {
  runQueue := make([]Runner, 0)
  f1: for {
    select {
    case runner := <- c.runChan:
      if c.runCount >= maxNumberRunningThreads {
        runQueue = append(runQueue, runner)
        log.Debugf("Queing %v...", len(runQueue))
      } else {
        c.runCount += 1
        c.startRunner(runner)
      }
    case <- c.runFinishedChan:
      if len(runQueue) > 0 {
        c.startRunner(runQueue[0])
        runQueue = runQueue[1:]
        log.Debugf("Running from queue %v...", len(runQueue))
      } else {
        c.runCount -= 1
      }

      if c.runCount == 0 {
        break f1
      }
    case backCount := <- c.backCountChan:
      if backCount {
        c.backCount += 1
      } else {
        c.backCount -= 1
      }
    case newError := <- c.errorChan:
      c.errors = append(c.errors, newError.Error())
    }
  }

  c.gracefullyClose()

  for c.backCount > 0 {
    backCount := <- c.backCountChan
    if backCount {
      c.backCount += 1
    } else {
      c.backCount -= 1
    }
  }

  c.finish()
}

func (c *Controller) gracefullyClose() {
  log.Debug("Gracefully Shuting down")
  close(c.doneChan)
}
func (c *Controller) finish() {
  log.Debug("Closing")
  close(c.finishedChan)
}

// will close when we want to safely shutdown
func (c *Controller) Done() chan struct{} {
  return c.doneChan
}
func (c *Controller) Wait() error {
  <- c.finishedChan
  if len(c.errors) == 0 {
    return nil
  }

  return fmt.Errorf(strings.Join(c.errors, "\n"))
}

func (c *Controller) Go(runner Runner) {
  c.runChan <- runner
}
func (c *Controller) startRunner(runner Runner) {
  go func() {
    runners, err := runner.Run(c)
    if err != nil {
      c.errorChan <- err
    } else {
      for _, r := range runners {
        c.Go(r)
      }
    }
    c.runFinishedChan <- true
  }()
}

// Running in background will call cancel even if still running
func (c *Controller) GoBackground(bgRunner BackgroundRunner) {
  c.backCountChan <- true
  go func() {
    err := bgRunner.Run(c)
    if err != nil {
      c.errorChan <- err
      c.gracefullyClose()
    }
    c.backCountChan <- false
  }()
}
