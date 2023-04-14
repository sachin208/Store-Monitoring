package models

type TimezoneModel struct {
	StoreId     string `gorm:"store_id,omitempty"`
	TimezoneStr string `gorm:"timezone_str,omitempty"`
}
