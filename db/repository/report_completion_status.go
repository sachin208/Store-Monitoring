package repository

import (
	"log"
	"store_monitoring/db"
	"store_monitoring/models"
)

type ReportCompletionStatusRepository struct{}

func (r *ReportCompletionStatusRepository) MarkStatusOfReportInDb(reportId, status, reportPath string) error {
	dbConnection := db.GetDB()
	if dbConnection == nil {
		log.Printf("Unable to obtain DB in status repository")
		return nil
	}
	query := "INSERT INTO Store_Monitoring.Store_Completion_Status(report_id, status, csv_file_path) VALUES (?, ?, ?)"
	_, insertErr := dbConnection.Exec(query, reportId, status, reportPath)
	if insertErr != nil {
		log.Printf("Error while inserting data in report_completion_status repository , Error = " + insertErr.Error())
	}
	return insertErr
}

func (r *ReportCompletionStatusRepository) UpdateStatusOfReportByReportId(reportId, status string) error {
	dbConnection := db.GetDB()
	if dbConnection == nil {
		log.Printf("Unable to obtain DB in status repository")
		return nil
	}
	query := "UPDATE Store_Monitoring.Store_Completion_Status SET status = ? WHERE report_id = ?"
	_, updateErr := dbConnection.Exec(query, status, reportId)
	if updateErr != nil {
		log.Printf("Error while updating data in report_completion_status repository , Error = " + updateErr.Error())
	}
	return updateErr
}

func (r *ReportCompletionStatusRepository) GetCsvFilePathByReportId(reportId string) (models.ReportCompletionStatusModel, error) {
	dbConnection := db.GetDB()
	var response models.ReportCompletionStatusModel
	if dbConnection == nil {
		log.Printf("Unable to obtain DB in status repository")
		return response, nil
	}
	query := "SELECT * From Store_Monitoring.Store_Completion_Status WHERE report_id = ?"
	findErr := dbConnection.QueryRow(query, reportId).Scan(&response.ReportId, &response.Status, &response.CsvFilePath)
	if findErr != nil {
		log.Printf("Error while fetching data in report_completion_status repository , Error = " + findErr.Error())
	}
	return response, findErr
}
