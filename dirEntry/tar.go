package dirEntry

import (
	"archive/tar"
	"bufio"
	"io"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	log "github.com/sirupsen/logrus"
)

func (de *DirEntry) iterateTar(reader io.Reader) error {
	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()


		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("Error reading tar (%v)", err)
		}

		if header.Typeflag == tar.TypeReg {
			tarBufReader := bufio.NewReader(tarReader)

			fileBytes, err := tarBufReader.Peek(int(detectReadSize(header.Size)))
			if err != nil {
				return fmt.Errorf("Error reading bytes from tar for detect (%v)", err)
			}

			mtype := mimetype.Detect(fileBytes)
			if strings.HasPrefix(mtype.String(), "image") {
				err = de.addFile(filepath.Base(header.Name), tarBufReader, mtype)
			} else {
				log.Infof("Skipping %v not an image (%v)", header.Name, mtype.String())
			}

			if err != nil {
				return fmt.Errorf("Error adding file in tar %v", err)
			}
		}
	}

	return nil
}
