package repository

import (
	"log"
	"store_monitoring/db"
	"store_monitoring/models"
)

type BusinessHoursRepository struct{}

func (b *BusinessHoursRepository) InsertDataInDb(rows []models.BusinessHoursModel) error {
	dbConnection := db.GetDB()
	if dbConnection == nil {
		log.Printf("Unable to obtain DB in business_hours repository")
		return nil
	}
	for _, row := range rows {
		query := "INSERT INTO Store_Monitoring.Business_Hours (store_id, dayOfWeek, start_time_local, end_time_local) VALUES (?, ?, ?, ?)"
		_, insertErr := dbConnection.Exec(query, row.StoreId, row.DaysOfWeek, row.StartTimeLocal, row.EndTimeLocal)
		if insertErr != nil {
			log.Printf("Error while inserting data in business_hours repository , Error = " + insertErr.Error())
		}
	}
	return nil
}

func (b *BusinessHoursRepository) GetDataFromDb() ([]models.BusinessHoursModel, error) {
	dbConnection := db.GetDB()
	if dbConnection == nil {
		log.Printf("Unable to obtain DB coonection in business_hours repository")
		return nil, nil
	}
	var rows []models.BusinessHoursModel
	query := "SELECT * FROM Store_Monitoring.Business_Hours "
	response, err := dbConnection.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	for response.Next() {
		var row models.BusinessHoursModel
		response.Scan(&row.StoreId, &row.DaysOfWeek, &row.StartTimeLocal, &row.EndTimeLocal)
		rows = append(rows, row)
	}
	return rows, nil
}
