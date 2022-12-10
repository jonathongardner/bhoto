package fileInventory

import (
	"database/sql"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

// Primary key and unique key add automatic indexes so no need to index those fields
func (f *Fin) SetupDB() (error) {
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
	err = exec(db, `CREATE TABLE fileInfos (
		"sha1" CHARACTER(64) NOT NULL PRIMARY KEY,
		"filetype" VARCHAR(255) NOT NULL,
		"extension" VARCHAR(255) NOT NULL,
		"taken" INT NOT NULL
	);`)
	if err != nil {
		return fmt.Errorf("file infos table (%v)", err)
	}

	log.Info("Creating index file info table...")
	err = exec(db, "CREATE INDEX fileinfosTaken ON fileInfos(taken);")
	if err != nil {
		return fmt.Errorf("file infos table (%v)", err)
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
