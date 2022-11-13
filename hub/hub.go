package hub

import (
	"context"
	"fmt"
	"strings"

	"github.com/jonathongardner/bhoto/dirEntry"

	log "github.com/sirupsen/logrus"
)

type Hub struct {
	dirEntryIn        chan dirEntry.DirEntry
	dirEntryFinished  chan bool
	finished          chan error
	dirEntryProcessor dirEntry.DirEntryProcessor
}
func NewHub(initialDir string, dep dirEntry.DirEntryProcessor) *Hub {
	// make dirEntryIn 1 buffer so we can go ahead and add the initialDir
	h := &Hub{
		dirEntryIn: make(chan dirEntry.DirEntry, 1),
		dirEntryFinished: make(chan bool),
		finished: make(chan error),
		dirEntryProcessor: dep,
	}
	h.dirEntryIn <- dirEntry.DirEntry{Path: initialDir, IsDir: true,}
	return h
}

func (h *Hub) Process(ctx context.Context, maxNumberOfFileProcessors int) {
	defer close(h.finished)

	dirEntriesQueue := make([]dirEntry.DirEntry, 0)
	running := 0

	listen: for {
		select {
		case <- ctx.Done():
			break listen
		case newDE := <- h.dirEntryIn:
			if running >= maxNumberOfFileProcessors {
				dirEntriesQueue = append(dirEntriesQueue, newDE)
			} else {
				go h.runDEProceesor(newDE)
				running++
			}
		case <- h.dirEntryFinished:
			if len(dirEntriesQueue) > 0 {
				go h.runDEProceesor(dirEntriesQueue[0])
				dirEntriesQueue = dirEntriesQueue[1:]
			} else {
				running--
				if running == 0 && len(h.dirEntryIn) == 0 {
					break listen
				}
			}
		}
	}
	log.Infof("Closing hub (%v)...", running)

	for running > 0 {
		if running % 5 == 0 {
			log.Infof("Waiting on %v file processors", running)
		}
		select {
		case <- h.dirEntryIn:
			// do nothing, just catch so we dont block others
		case <- h.dirEntryFinished:
			running--
		}
	}
}

// Should only listen from one thread
func (h *Hub) Wait() (error) {
	errors := []string{}
	for {
		err, ok := <- h.finished
		if !ok {
			break
		}
		errors = append(errors, err.Error())
	}
	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "\n"))
	}
	return nil
}

func (h *Hub) runDEProceesor(de dirEntry.DirEntry) {
	dirEntriesToAdd, err := h.dirEntryProcessor.ProcessDirEntry(de)
	if err != nil {
		h.finished <- err
	}

	for _, deNew := range dirEntriesToAdd {
		h.dirEntryIn <- deNew
	}

	h.dirEntryFinished <- true
}
