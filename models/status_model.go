package models

import "time"

type StatusModel struct {
	StoreId      string    `gorm:"store_id,omitempty"`
	TimeStampUtc time.Time `gorm:"timestamp_utc,omitempty"`
	Status       string    `gorm:"status,omitempty"`
}
