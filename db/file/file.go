package file

import (
	"encoding/csv"
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"store_monitoring/constants"
	"store_monitoring/db/repository"
	"store_monitoring/dto"
	"store_monitoring/models"
	"strconv"
	"time"
)

type File struct{}

func (f *File) PrepareDatabase() {
	// Updating data in to DB from all CSV Files
	f.ReadCsvFilesAndUpdateDB()
}

func (f *File) ReadCsvFilesAndUpdateDB() {
	f.ReadStatusFileAndUpdateDb(constants.StatusFilePath)
	f.ReadBusinessHoursFileAndUpdateDb(constants.BusinessHoursFilePath)
	f.ReadTimeZoneFileAndUpdateDb(constants.TimezoneFilePath)
}

func (f *File) PrepareCsvForReport(response map[string]dto.UpAndDownTimeResponseDto) string {
	reportId := uuid.New().String()
	fileName := reportId + "_report_file"
	filePath := constants.FileFolderPath + fileName
	reportCompletionStatusRepository := repository.ReportCompletionStatusRepository{}
	reportStatus := "PENDING"
	err := reportCompletionStatusRepository.MarkStatusOfReportInDb(reportId, reportStatus, filePath)
	if err != nil {
		fmt.Print("Error while updating status in report_completion_status, Error = ", err.Error())
	}
	go f.CreateReportCsvFile(filePath, response, reportStatus, reportId)

	return reportId
}

func (f *File) CreateReportCsvFile(filePath string, storeDataMap map[string]dto.UpAndDownTimeResponseDto, reportStatus, reportId string) {
	csvFile, err := os.Create(filePath)
	if err != nil {
		fmt.Print("Error while creating new file, Error = ", err.Error())
		return
	}
	csvWriter := csv.NewWriter(csvFile)
	rows := [][]string{
		{"store_id", "uptime_last_hour", "uptime_last_day", "update_last_week", "downtime_last_hour", "downtime_last_day", "downtime_last_week"},
	}
	for storeId, storeData := range storeDataMap {
		row := []string{storeId, storeData.UpTimeLastHour, storeData.UpTimeLastDay, storeData.UpTimeLastWeek, storeData.DownTimeLastHour, storeData.DownTimeLastDay, storeData.DownTimeLastWeek}
		rows = append(rows, row)
	}
	err = csvWriter.WriteAll(rows)

	// Update status in DB
	if err != nil {
		reportStatus = "FAILED"
	} else {
		reportStatus = "COMPLETED"
	}
	reportCompletionStatusRepository := repository.ReportCompletionStatusRepository{}
	err = reportCompletionStatusRepository.UpdateStatusOfReportByReportId(reportId, reportStatus)
	if err != nil {
		fmt.Print("Error while updating status in report_completion_status, Error = ", err.Error())
	}
}

func (f *File) ReadStatusFileAndUpdateDb(filePath string) {
	rows := f.ReadFile(filePath)
	var data []models.StatusModel
	for _, row := range rows {
		if row[0] == "store_id" {
			continue
		}
		storeId := row[0]
		status := row[1]
		timeStamp := ParseTimeFromString(row[2])
		currentRow := models.StatusModel{
			StoreId:      storeId,
			TimeStampUtc: timeStamp.UTC(),
			Status:       status,
		}
		data = append(data, currentRow)
	}
	statusRepository := repository.StatusRepository{}
	err := statusRepository.InsertDataInDb(data)
	if err != nil {
		log.Printf("Error while inserting data in Status table, Error = " + err.Error())
	}
}

func (f *File) ReadBusinessHoursFileAndUpdateDb(filePath string) {
	rows := f.ReadFile(filePath)
	var data []models.BusinessHoursModel
	for _, row := range rows {
		if row[0] == "store_id" {
			continue
		}
		storeId := row[0]
		if row[1] == "" {
			row[1] = "7"
		}
		daysOfWeek, _ := strconv.Atoi(row[1])
		startTimeLocal := row[2]
		endTimeLocal := row[3]
		currentRow := models.BusinessHoursModel{
			StoreId:        storeId,
			DaysOfWeek:     daysOfWeek,
			StartTimeLocal: startTimeLocal,
			EndTimeLocal:   endTimeLocal,
		}
		data = append(data, currentRow)
	}
	businessHoursRepository := repository.BusinessHoursRepository{}
	err := businessHoursRepository.InsertDataInDb(data)
	if err != nil {
		log.Printf("Error while inserting data in business_hours collection, Error = " + err.Error())
	}
}

func (f *File) ReadTimeZoneFileAndUpdateDb(filePath string) {
	rows := f.ReadFile(filePath)
	var data []models.TimezoneModel
	defaultTimezone := "America/Chicago"
	for _, row := range rows {
		if row[0] == "store_id" {
			continue
		}
		storeId := row[0]
		timeZoneStr := row[1]
		if timeZoneStr == "" {
			timeZoneStr = defaultTimezone
		}
		currentRow := models.TimezoneModel{
			StoreId:     storeId,
			TimezoneStr: timeZoneStr,
		}
		data = append(data, currentRow)
	}
	timezoneRepository := repository.TimezoneRepository{}
	err := timezoneRepository.InsertDataInDb(data)
	if err != nil {
		log.Printf("Error while inserting data in timezone collection, Error = " + err.Error())
	}
}

func (f *File) ReadFile(filePath string) [][]string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(file)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func ParseTimeFromString(timeString string) time.Time {
	stringLength := len(timeString)
	newTimeString := timeString[0:10] + "T" + timeString[11:(stringLength-4)] + "Z"
	t, _ := time.Parse(time.RFC3339, newTimeString)
	return t
}
