package fileInventory

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type dbRunner interface {
	Run(db *sql.DB) error
}

func SetupDB(path string) (error) {
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("database already exist (%v)", path)
	}

	log.Infof("Creating database at %v", path)
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error opneing database (%v - %v)", path, err)
	}
	file.Close()

	log.Info("Opening database")
	db, _ := sql.Open("sqlite3", path)
	defer db.Close()

	//---------------FileInfo----------------
	log.Info("Creating file info table...")
	createFileInfoTableSQL := `CREATE TABLE fileInfos (
		"sha1" CHARACTER(64) NOT NULL PRIMARY KEY,
		"filetype" VARCHAR(255) NOT NULL,
		"extension" VARCHAR(255) NOT NULL
	);`
	statement1, err := db.Prepare(createFileInfoTableSQL) // Prepare SQL Statement
	if err != nil {
		return fmt.Errorf("error preparing file infos table (%v)", err)
	}
	_, err = statement1.Exec() // Execute SQL Statements
	if err != nil {
		return fmt.Errorf("error creating file infos table (%v)", err)
	}
	//---------------FileInfo----------------

	//---------------File----------------
	log.Info("Creating file table...")
	createFileTableSQL := `CREATE TABLE files (
		"sha1" CHARACTER(64) NOT NULL PRIMARY KEY,
		"file" BLOB NOT NULL
	);`
	statement2, err := db.Prepare(createFileTableSQL) // Prepare SQL Statement
	if err != nil {
		return fmt.Errorf("error preparing files table (%v)", err)
	}
	_, err = statement2.Exec() // Execute SQL Statements
	if err != nil {
		return fmt.Errorf("error creating files table (%v)", err)
	}
	//---------------File----------------

	//---------------FileNames----------------
	log.Info("Creating file name table...")
	createFileNameTableSQL := `CREATE TABLE fileNames (
		"id" INT AUTO INCREMENT PRIMARY KEY,
		"sha1" CHARACTER(64) NOT NULL,
		"name" VARCHAR(255) NOT NULL,
		UNIQUE(sha1, name)
	);`
	statement3, err := db.Prepare(createFileNameTableSQL) // Prepare SQL Statement
	if err != nil {
		return fmt.Errorf("error preparing file names table (%v)", err)
	}
	_, err = statement3.Exec() // Execute SQL Statements
	if err != nil {
		return fmt.Errorf("error creating file names table (%v)", err)
	}
	//---------------FileNames----------------

	log.Info("Created file name table")

	return nil
}

func (f *Fin) StartDB(path string) {
	if _, err := os.Stat(path); err != nil {
		log.Errorf("database doesn't exist (%v)", path)
		return
	}

	db, _ := sql.Open("sqlite3", path)
	defer db.Close()

	count := uint64(0)
	for {
		dbRunner, ok := <- f.db
		if !ok {
			log.Infof("Processed %v", count)
			break
		}

		err := dbRunner.Run(db)
		if err != nil {
			log.Errorf("Error adding to db %v", err)
		}
		count += 1
		if count % 100 == 0 {
			log.Infof("Processed %v", count)
		}
	}
}

type file struct {
  filename  string
	sha1      string
  filetype  string
  extension string
  fileBytes []byte
}

func (f *Fin) addFile(filename string, sha1 string, filetype string, extension string, fileBytes []byte) {
	f.db <- &file{filename: filename, sha1: sha1, filetype: filetype, extension: extension, fileBytes: fileBytes, }
}

func (fi *file) Run(db *sql.DB) error {
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
