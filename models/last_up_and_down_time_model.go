package models

type LastUpAndDownTimeModel struct {
	StoreId          string `bson:"store_id,omitempty"`
	UpTimeLastHour   string `bson:"uptime_last_hour,omitempty"`
	UpTimeLastDay    string `bson:"uptime_last_day,omitempty"`
	UpTimeLastWeek   string `bson:"update_last_week,omitempty"`
	DownTimeLastHour string `bson:"downtime_last_hour,omitempty"`
	DownTimeLastDay  string `bson:"downtime_last_day,omitempty"`
	DownTimeLastWeek string `bson:"downtime_last_week,omitempty"`
}
