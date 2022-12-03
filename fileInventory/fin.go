package fileInventory

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jonathongardner/bemery/routines"

	log "github.com/sirupsen/logrus"
)

type Fin struct {
	ImageDB   chan dbRunner
	VideoDB   chan dbRunner
}

func NewFin(path string) (*Fin, error) {
	return &Fin{ImageDB: make(chan dbRunner), ImageDB: make(chan dbRunner)}, nil
}

type DB struct {
	path  string
	input chan dbRunner
}

type dbRunner interface {
	Run(db *sql.DB) error
}

func NewDB(path string, input chan dbRunner) (*DB, error) {
	return &Fin{input: input, path: path}, nil
}

func (f *DB) dbExist() bool {
	_, err := os.Stat(f.path)
	return err == nil
}

// only run one cause dont want mulitple sqlite dbs open
func (f *DB) Run(rc *routines.Controller) error {
	if !f.dbExist() {
		return fmt.Errorf("database doesn't exist (%v)", f.path)
	}

	db, _ := sql.Open("sqlite3", f.path)
	defer db.Close()

	count := uint64(0)
	f1: for {
		select {
		case dbRunner := <- f.input:
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

// https://zetcode.com/golang/sqlite3/
func (f *DB) Stats() error {
	if !f.dbExist() {
		return fmt.Errorf("database doesn't exist (%v)", f.path)
	}

	db, _ := sql.Open("sqlite3", f.path)
	defer db.Close()

	fileNameStats := `SELECT count(distinct sha1) as uniqueFiles, count(name) as total FROM fileNames;`
	statement1, err := db.Prepare(fileNameStats) // Prepare SQL Statement
	if err != nil {
		return fmt.Errorf("error getting filename stats (%v)", err)
	}
	var uniqueFiles string
	var name string
	err = statement1.QueryRow().Scan(&uniqueFiles, &name)
	if err != nil {
		return fmt.Errorf("error getting filename stats (%v)", err)
	}
	log.Infof("Unique Files: %v, Total: %v", uniqueFiles, name)

	return nil
}

// Primary key and unique key add automatic indexes so no need to index those fields
func (f *DB) SetupDB() (error) {
	if f.dbExist() {
		return fmt.Errorf("database already exist (%v)", f.path)
	}

	log.Infof("Creating database at %v", f.path)
	file, err := os.Create(f.path)
	if err != nil {
		return fmt.Errorf("error opneing database (%v - %v)", f.path, err)
	}
	file.Close()

	log.Info("Opening database")
	db, _ := sql.Open("sqlite3", f.path)
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
