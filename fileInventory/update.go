package fileInventory

import (
	"database/sql"
	"fmt"
	"errors"

	"github.com/jonathongardner/bhoto/photo"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

// end of data
var eod = errors.New("EOD")
const updatePaging = 10

func (f *Fin) RebuildFileInfo() error {
	if !f.dbExist() {
		return fmt.Errorf("database doesn't exist (%v)", f.path)
	}

	db, _ := sql.Open("sqlite3", f.path)
	defer db.Close()

	previosSha256 := ""
	count := 0
	for {
		imgs, err := nextFileRows(db, previosSha256)

		for _, img := range imgs {
			previosSha256 = img.Sha256()
			img.UdateFileInto(db)
			count = count + 1
		}

		if err == eod {
			break
		}
		if err != nil {
			return err
		}

		if count % 1000 == 0 {
			log.Infof("Finished %v", count)
		}
	}

	log.Infof("Finished %v", count)
	return nil
}

func nextFileRows(db *sql.DB, sha256 string) ([]*photo.Image, error) {
	rows, err := db.Query(`SELECT sha1, file FROM files WHERE sha1 > ? ORDER BY sha1 ASC LIMIT ?;`, sha256, updatePaging)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	toReturn := []*photo.Image{}
	for rows.Next() {
		sha256 := ""
		fb := []byte{}
		rows.Scan(&sha256, &fb)
		img := photo.NewImageWithChecksum("", fb, sha256)

		img.SetMagicInfo()
		img.SetExifInfo()

		toReturn = append(toReturn, img)
	}
	if len(toReturn) < updatePaging {
		return toReturn, eod
	}
	return toReturn, nil
}
