package fileInventory

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

  "github.com/jonathongardner/bhoto/dirEntry"

	"github.com/gabriel-vasile/mimetype"
	log "github.com/sirupsen/logrus"
)

// 100 MB
// TODO: Config
const maxFileSize int64 = 100 * 1024 *1024

type Fin struct {
	db chan dbRunner
}

func NewFin() (*Fin, error) {
	return &Fin{db: make(chan dbRunner, 5)}, nil
}

func (f *Fin) Close() error {
	close(f.db)
	return nil
}

// return a list of DirEntry's to add if any
func (fin *Fin) ProcessDirEntry(de dirEntry.DirEntry) ([]dirEntry.DirEntry, error) {
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

	fin.addFile(filename, checksum, mtype.String(), mtype.Extension(), fileBytes)

  return []dirEntry.DirEntry{}, nil
}
