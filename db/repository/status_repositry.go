package repository

import (
	"log"
	"store_monitoring/db"
	"store_monitoring/models"
	"time"
)

type StatusRepository struct{}

func (s *StatusRepository) InsertDataInDb(rows []models.StatusModel) error {
	dbConnection := db.GetDB()
	if dbConnection == nil {
		log.Printf("Unable to obtain DB in status repository")
		return nil
	}
	for _, row := range rows {
		query := "INSERT INTO Store_Monitoring.Status(store_id, timestamp_utc, status) VALUES (?, ?, ?)"
		_, insertErr := dbConnection.Exec(query, row.StoreId, row.TimeStampUtc, row.Status)
		if insertErr != nil {
			log.Printf("Error while inserting data in status repository , Error = " + insertErr.Error())
		}
	}
	return nil
}

func (s *StatusRepository) GetDataFromDb() ([]models.StatusModel, error) {
	dbConnection := db.GetDB()
	if dbConnection == nil {
		log.Printf("Unable to obtain DB coonection in status repository")
		return nil, nil
	}
	var rows []models.StatusModel
	lastWeekTimeStamp := time.Now().AddDate(0, 0, -7*100).UTC().String()
	query := "SELECT * FROM Store_Monitoring.Status WHERE timestamp_utc >= ? ORDER by timestamp_utc ASC"
	response, err := dbConnection.Query(query, lastWeekTimeStamp)
	if err != nil {
		log.Fatal(err)
	}
	for response.Next() {
		var storeId string
		var status string
		var timeStamp []uint8
		err = response.Scan(&storeId, &timeStamp, &status)
		t, _ := time.Parse(time.RFC3339, string(timeStamp))
		rows = append(rows, models.StatusModel{
			StoreId:      storeId,
			Status:       status,
			TimeStampUtc: t,
		})
	}
	return rows, nil
}
