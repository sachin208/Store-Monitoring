package models

type BusinessHoursModel struct {
	StoreId        string `gorm:"store_id,omitempty"`
	DaysOfWeek     int    `gorm:"dayOfWeek,omitempty"`
	StartTimeLocal string `gorm:"start_time_local,omitempty"`
	EndTimeLocal   string `gorm:"end_time_local,omitempty"`
}
