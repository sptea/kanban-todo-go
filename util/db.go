package util

import (
	"database/sql"
	"log"
)

var Db *sql.DB

func InitDB() {
	var err error
	if Db, err = sql.Open("sqlite3", DbPath); err != nil {
		log.Printf("Couldnt open database file. FIlePath: " + DbPath)
		log.Panic(err)
	}

	if err = Db.Ping(); err != nil {
		log.Panic(err)
	}
}
