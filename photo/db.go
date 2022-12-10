package photo

import (
  "database/sql"
	// "time"

	// log "github.com/sirupsen/logrus"
)

func (fi *Image) UdateFileInto(db *sql.DB) error {
  insertFileInfoSQL := `UPDATE fileInfos SET filetype = ?, extension = ?, taken = ? WHERE sha1 = ?`
  // insertFileInfoSQL := `INSERT OR REPLACE INTO fileInfos(sha1, type) VALUES (?, ?)`
  statement1, err := db.Prepare(insertFileInfoSQL)
  if err != nil {
    return err
  }

  _, err = statement1.Exec(fi.filetype, fi.extension, fi.taken, fi.sha256)
  return err
}

func (fi *Image) Run(db *sql.DB) error {
	// time.Sleep(8 * time.Second)
  //---------------FileInfo----------------
  insertFileInfoSQL := `INSERT OR IGNORE INTO fileInfos(sha1, filetype, extension, taken) VALUES (?, ?, ?, ?)`
	// insertFileInfoSQL := `INSERT OR REPLACE INTO fileInfos(sha1, type) VALUES (?, ?)`
	statement1, err := db.Prepare(insertFileInfoSQL)
	if err != nil {
		return err
	}

	_, err = statement1.Exec(fi.sha256, fi.filetype, fi.extension, fi.taken)
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

	_, err = statement2.Exec(fi.sha256, fi.fileBytes)
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

	_, err = statement3.Exec(fi.sha256, fi.filename)
	if err != nil {
		return err
	}
	return nil
  //---------------FileName----------------
}
