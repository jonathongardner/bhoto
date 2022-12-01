package fileInventory

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jonathongardner/bhoto/routines"

	log "github.com/sirupsen/logrus"
)

type Fin struct {
	path string
	db   chan dbRunner
}

type dbRunner interface {
	Run(db *sql.DB) error
}

func NewFin(path string) (*Fin, error) {
	return &Fin{db: make(chan dbRunner, 5), path: path}, nil
}

func (f *Fin) dbExist() bool {
	_, err := os.Stat(f.path)
	return err == nil
}

// only run one cause dont want mulitple sqlite dbs open
func (f *Fin) Run(rc *routines.Controller) error {
	if !f.dbExist() {
		return fmt.Errorf("database doesn't exist (%v)", f.path)
	}

	db, _ := sql.Open("sqlite3", f.path)
	defer db.Close()

	count := uint64(0)
	f1: for {
		select {
		case dbRunner := <- f.db:
			err := dbRunner.Run(db)
			if err != nil {
				log.Errorf("Error adding to db %v", err)
			}
			count += 1
			if count % 100 == 0 {
				log.Infof("Processed %v", count)
			}
		case <- rc.Done():
			log.Infof("Processed %v", count)
			break f1
		}
	}
	return nil
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
