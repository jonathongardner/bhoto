package dirEntry

import (
	"compress/gzip"
	"io"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jonathongardner/bhoto/photo"
	"github.com/jonathongardner/bhoto/fileInventory"
	"github.com/jonathongardner/bhoto/routines"

	"github.com/gabriel-vasile/mimetype"
	log "github.com/sirupsen/logrus"
)

// https://github.com/gabriel-vasile/mimetype/blob/master/mimetype.go#L17
const maxBytesFileDetect int64 = 3072

func detectReadSize(s int64) int64 {
	if maxBytesFileDetect > s {
		return s
	}
	return maxBytesFileDetect
}

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

	file, err := os.Open(de.Path)
	if err != nil {
		return nil, fmt.Errorf("Error opening directory entry (%v)", err)
	}
	defer file.Close()

	// need to read atleast 512 for some reaons
	fileBytes := make([]byte, detectReadSize(de.Size))
	_, err = file.Read(fileBytes)
	if err != nil {
		return nil, fmt.Errorf("Error reading bytes (%v)", err)
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("Error rewind (%v)", err)
	}

	mtype := mimetype.Detect(fileBytes)

	added, err := de.addFile(de.Path, file, mtype)
	if !added {
		if mtype.String() == "application/zip" {
			log.Debugf("Processing %v as zip", de.Path)
			err = de.iterateZip(file, de.Size)
		} else if mtype.String() == "application/x-tar" {
			log.Debugf("Processing %v as tar", de.Path)
			err = de.iterateTar(file)
		} else if mtype.String() == "application/gzip" {
			log.Debugf("Processing %v as gtar", de.Path)
			gzf, err := gzip.NewReader(file)
			if err != nil {
				return nil, fmt.Errorf("Error opening gzip (%v)", err)
			}
			err = de.iterateTar(gzf)
		}
	}

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (de *DirEntry) addFile(path string, reader io.Reader, mtype *mimetype.MIME) (bool, error) {
	mgroup := strings.SplitN(mtype.String(), "/", 2)[0]
	// TODO: make config
	if mgroup != "image" && mgroup != "video" {
		log.Infof("Skipping %v not an image or video (%v - %v)", path, mgroup, mtype.String())
		return false, nil
	}
	fileBytes, err := io.ReadAll(reader)
	if err != nil {
		return true, fmt.Errorf("Error reading rest of bytes (%v)", err)
	}

	img := photo.NewImageWithMagic(filepath.Base(path), fileBytes, mtype)
	img.SetChecksum()
	img.SetExifInfo()

	de.fin.AddFile(img)

  return true, nil
}
