package dirEntry

import (
	"fmt"
	"os"
	"path/filepath"
)

type DirEntry struct {
	Path  string
	IsDir bool
	Size  int64
}

func (de DirEntry) Children() ([]DirEntry, error) {
	toReturn := make([]DirEntry, 0)

	dirEntries, err := os.ReadDir(de.Path)
	if err != nil {
		return toReturn, fmt.Errorf("Error opening dir %v", err)
	}

	for _, newDE := range dirEntries {
		path := filepath.Join(de.Path, newDE.Name())
		if newDE.IsDir() {
			toReturn = append(toReturn, DirEntry{ Path: path, IsDir: true, Size: -1, })
		} else {
			info, err := newDE.Info()
			if err != nil {
				return nil, fmt.Errorf("Error getting dir info %v", err)
			}
			toReturn = append(toReturn, DirEntry{ Path: path, IsDir: false, Size: info.Size(), })
		}
	}
	return toReturn, nil
}

type DirEntryProcessor interface {
	ProcessDirEntry(DirEntry) ([]DirEntry, error)
}
