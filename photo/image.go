package photo

import (
  "crypto/sha256"
  "encoding/hex"
	"fmt"
	"time"

  "github.com/gabriel-vasile/mimetype"

	// log "github.com/sirupsen/logrus"
)
var noTime = time.Date(1970, time.Month(1), 1, 0, 0, 0, 0, time.UTC)

type Image struct {
  filename  string
	sha256      string
  filetype  string
	extension string
  taken time.Time
  fileBytes []byte
}

func (fi *Image) Sha256() string {
  return fi.sha256
}

func NewImage(filename string, fileBytes []byte) (*Image) {
	return &Image{filename: filename, fileBytes: fileBytes, taken: noTime, }
}

func NewImageWithMagic(filename string, fileBytes []byte, mtype *mimetype.MIME) *Image {
	toReturn := NewImage(filename, fileBytes)
  toReturn.updateMagic(mtype)
  return toReturn
}
func (fi *Image) updateMagic(mtype *mimetype.MIME) {
  fi.filetype = mtype.String()
  fi.extension = mtype.Extension()
}
func (fi *Image) SetMagicInfo() error {
  if fi.filetype != "" && fi.extension != "" {
    return fmt.Errorf("filetype and extension already set")
  }

  fi.updateMagic(mimetype.Detect(fi.fileBytes))
  return nil
}

func NewImageWithChecksum(filename string, fileBytes []byte, sha256 string) *Image {
	toReturn := NewImage(filename, fileBytes)
  toReturn.sha256 = sha256
  return toReturn
}
func (fi *Image) SetChecksum() error {
  if fi.sha256 != "" {
    return fmt.Errorf("sha256 already set")
  }

  hash := sha256.Sum256(fi.fileBytes)
  fi.sha256 = hex.EncodeToString(hash[:])
  return nil
}
