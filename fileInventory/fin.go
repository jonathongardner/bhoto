package fileInventory

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jonathongardner/bhoto/photo"
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
	return &Fin{db: make(chan dbRunner), path: path}, nil
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

func (f *Fin) AddFile(img *photo.Image) {
	f.db <- img
}

// https://zetcode.com/golang/sqlite3/
func (f *Fin) Stats() error {
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

func exec(db *sql.DB, q string) error {
	statement1, err := db.Prepare(q) // Prepare SQL Statement
	if err != nil {
		return fmt.Errorf("error preparing %v", err)
	}
	_, err = statement1.Exec() // Execute SQL Statements
	if err != nil {
		return fmt.Errorf("error running %v", err)
	}
	return nil
}
