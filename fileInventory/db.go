package fileInventory

import (
	"database/sql"
	// "time"

	_ "github.com/mattn/go-sqlite3"
	// log "github.com/sirupsen/logrus"
)

type image struct {
  filename  string
	sha1      string
  filetype  string
  extension string
  fileBytes []byte
}

func (f *Fin) AddFile(filename string, sha1 string, filetype string, extension string, fileBytes []byte) {
	f.db <- &image{filename: filename, sha1: sha1, filetype: filetype, extension: extension, fileBytes: fileBytes, }
}

func (fi *image) Run(db *sql.DB) error {
	// time.Sleep(8 * time.Second)
  //---------------FileInfo----------------
  insertFileInfoSQL := `INSERT OR IGNORE INTO fileInfos(sha1, filetype, extension) VALUES (?, ?, ?)`
	// insertFileInfoSQL := `INSERT OR REPLACE INTO fileInfos(sha1, type) VALUES (?, ?)`
	statement1, err := db.Prepare(insertFileInfoSQL)
	if err != nil {
		return err
	}

	_, err = statement1.Exec(fi.sha1, fi.extension, fi.extension)
	if err != nil {
		return err
	}
  //---------------FileInfo----------------

	//---------------File----------------
	insertFileSQL := `INSERT OR IGNORE INTO files(sha1, file) VALUES (?, ?)`
	statement2, err := db.Prepare(insertFileSQL)
	if err != nil {
		return err
	}

	_, err = statement2.Exec(fi.sha1, fi.fileBytes)
	if err != nil {
		return err
	}
  //---------------File----------------

	//---------------FileName----------------
	insertFileNameSQL := `INSERT OR IGNORE INTO fileNames(sha1, name) VALUES (?, ?)`
	statement3, err := db.Prepare(insertFileNameSQL)
	if err != nil {
		return err
	}

	_, err = statement3.Exec(fi.sha1, fi.filename)
	if err != nil {
		return err
	}
	return nil
  //---------------FileName----------------
}
