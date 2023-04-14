package dto

type UpAndDownTimeResponseDto struct {
	UpTimeLastHour   string // in minutes
	UpTimeLastDay    string // in hours
	UpTimeLastWeek   string // in hours
	DownTimeLastHour string
	DownTimeLastDay  string
	DownTimeLastWeek string
}
