package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"store_monitoring/constants"
)

var dbService *sql.DB

func InitDB() {
	db, err := sql.Open("mysql", constants.DbConnectionString)
	if err != nil {
		log.Println("Error while connecting DB", err)
		panic(err)
	}
	dbService = db
}

func GetDB() *sql.DB {
	return dbService
}
