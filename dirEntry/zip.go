package dirEntry

import (
	"archive/zip"
	"bufio"
	"io"
	"fmt"

	"github.com/gabriel-vasile/mimetype"
	// log "github.com/sirupsen/logrus"
)

func (de *DirEntry) iterateZip(reader io.ReaderAt, size int64) error {
	zipReader, err := zip.NewReader(reader, size)
	if err == io.EOF {
		return fmt.Errorf("Error opening zip (%v)", err)
	}

	for _, file := range zipReader.File {
		fileInfo := file.FileInfo()

		if fileInfo.IsDir() {
			continue
		}

		fileInArchive, err := file.Open()
		if err != nil {
			return fmt.Errorf("Error opening archive file %v (%v)", file.Name, err)
		}

		tarBufReader := bufio.NewReader(fileInArchive)

		fileBytes, err := tarBufReader.Peek(int(detectReadSize(fileInfo.Size())))
		if err != nil {
			fileInArchive.Close()
			return fmt.Errorf("Error reading bytes from zip for detect (%v)", err)
		}

		mtype := mimetype.Detect(fileBytes)
		_, err = de.addFile(file.Name, tarBufReader, mtype)

		fileInArchive.Close()

		if err != nil {
			return fmt.Errorf("Error adding file in tar %v", err)
		}
	}

	return nil
}
