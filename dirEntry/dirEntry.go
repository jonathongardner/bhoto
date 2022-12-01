package dirEntry

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jonathongardner/bhoto/fileInventory"
	"github.com/jonathongardner/bhoto/routines"

	"github.com/gabriel-vasile/mimetype"
	log "github.com/sirupsen/logrus"
)

// 100 MB
// TODO: Config
const maxFileSize int64 = 100 * 1024 *1024

type DirEntry struct {
	Path  string
	IsDir bool
	Size  int64
	fin   *fileInventory.Fin
}

func NewDirEntry(path string, fin *fileInventory.Fin) (*DirEntry, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return &DirEntry{ Path: path, IsDir: info.IsDir(), Size: info.Size(), fin: fin, }, nil
}

func (de *DirEntry) Children() ([]routines.Runner, error) {
	toReturn := make([]routines.Runner, 0)

	dirEntries, err := os.ReadDir(de.Path)
	if err != nil {
		return toReturn, fmt.Errorf("Error opening dir %v", err)
	}

	for _, newDE := range dirEntries {
		path := filepath.Join(de.Path, newDE.Name())
		if newDE.IsDir() {
			toReturn = append(toReturn, &DirEntry{ Path: path, IsDir: true, Size: -1, fin: de.fin, })
		} else {
			info, err := newDE.Info()
			if err != nil {
				return nil, fmt.Errorf("Error getting dir info %v", err)
			}
			toReturn = append(toReturn, &DirEntry{ Path: path, IsDir: false, Size: info.Size(), fin: de.fin, })
		}
	}
	return toReturn, nil
}

// return a list of controller.Runner (if any) to run
func (de *DirEntry) Run(rc *routines.Controller) ([]routines.Runner, error) {
  if de.IsDir {
    return de.Children()
  }

	if de.Size > maxFileSize {
		log.Infof("Skipping %v larging then max size (%v - %v)", de.Path, de.Size, maxFileSize)
		return nil, nil
	}

	fileBytes, err := os.ReadFile(de.Path)
	if err != nil {
		return nil, fmt.Errorf("Error opening directory entry (%v)", err)
	}

	mtype := mimetype.Detect(fileBytes)
	if !strings.HasPrefix(mtype.String(), "image") {
		log.Infof("Skipping %v not an image (%v)", de.Path, mtype.String())
		return nil, nil
	}

	filename := filepath.Base(de.Path)

	hash := sha256.Sum256(fileBytes)
  checksum := hex.EncodeToString(hash[:])

	de.fin.AddFile(filename, checksum, mtype.String(), mtype.Extension(), fileBytes)

  return nil, nil
}
