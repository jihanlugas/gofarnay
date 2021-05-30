package config

import "database/sql"

func DbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "jihanlugas"
	dbPass := "12345678"
	dbName := "gofarnay"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName+"?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	return db
}
