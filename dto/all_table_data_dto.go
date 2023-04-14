package dto

import "time"

type AllTableDataDto struct {
	StatusWithTimeList              []StatusWithTimeDto
	WeekDaysWithStartAndEndTimeList []WeekDaysWithStartAndEndTime
	Timezone                        string
}

type StatusWithTimeDto struct {
	TimestampUtc time.Time
	Status       string
}

type WeekDaysWithStartAndEndTime struct {
	DayOfWeek int
	StartTime string
	EndTime   string
}
