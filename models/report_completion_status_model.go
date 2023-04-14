package models

type ReportCompletionStatusModel struct {
	ReportId    string `gorm:"report_id,omitempty"`
	Status      string `gorm:"status,omitempty"`
	CsvFilePath string `gorm:"csv_file_path,omitempty"`
}
