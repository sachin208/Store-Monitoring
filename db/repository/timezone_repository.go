package repository

import (
	"log"
	"store_monitoring/db"
	"store_monitoring/models"
)

type TimezoneRepository struct{}

func (t *TimezoneRepository) InsertDataInDb(rows []models.TimezoneModel) error {
	dbConnection := db.GetDB()
	if dbConnection == nil {
		log.Printf("Unable to obtain DB in timezone repository")
		return nil
	}
	for _, row := range rows {
		query := "INSERT INTO Store_Monitoring.Timezone (store_id, timezone_str) VALUES (?, ?)"
		_, insertErr := dbConnection.Exec(query, row.StoreId, row.TimezoneStr)
		if insertErr != nil {
			log.Printf("Error while inserting data in timezone repository , Error = " + insertErr.Error())
		}
	}
	return nil
}

func (t *TimezoneRepository) GetDataFromDb() ([]models.TimezoneModel, error) {
	dbConnection := db.GetDB()
	if dbConnection == nil {
		log.Printf("Unable to obtain DB coonection in timezone repository")
		return nil, nil
	}
	var rows []models.TimezoneModel
	query := "SELECT * FROM Store_Monitoring.Timezone"
	response, err := dbConnection.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	for response.Next() {
		var row models.TimezoneModel
		response.Scan(&row.StoreId, &row.TimezoneStr)
		rows = append(rows, row)
	}
	return rows, nil
}
