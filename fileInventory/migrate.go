package fileInventory

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Primary key and unique key add automatic indexes so no need to index those fields
func (f *Fin) MigrateDB() (error) {
	if !f.dbExist() {
		return fmt.Errorf("database doesn't exist (%v)", f.path)
	}

	db, _ := sql.Open("sqlite3", f.path)
	defer db.Close()

	//---------------FileInfo----------------
	log.Info("Creating file info table...")
	err := exec(db, `ALTER TABLE fileInfos ADD COLUMN taken INT NOT NULL DEFAULT 0;`)
	if err != nil {
		return fmt.Errorf("file infos table (%v)", err)
	}

	log.Info("Creating index file info table...")
	err = exec(db, "CREATE INDEX fileinfosTaken ON fileInfos(taken);")
	if err != nil {
		return fmt.Errorf("file infos table (%v)", err)
	}
	//---------------FileInfo----------------

	return nil
}
