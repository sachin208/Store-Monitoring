package main

import (
	"store_monitoring/db"
	"store_monitoring/db/file"
	"store_monitoring/router"
)

func main() {
	db.InitDB()
	fileInstance := file.File{}
	fileInstance.PrepareDatabase()
	router.InitRouter()
}
