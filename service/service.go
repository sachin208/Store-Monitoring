package service

import (
	"fmt"
	"store_monitoring/constants"
	"store_monitoring/db/repository"
	"store_monitoring/dto"
	"strconv"
	"time"
)

type Service struct{}

var storeIdToAllTableDataMap = make(map[string]dto.AllTableDataDto)

func (s *Service) CalculateLastUpAndDownTime() map[string]dto.UpAndDownTimeResponseDto {
	s.appendStatusTableDataInMap()
	s.appendBusinessHoursTableDataInMap()
	s.appendTimezoneTableDataInMap()
	lastUpAndDownTimeResponseDtoMap := s.performBusinessLogicOfCalculatingUpAndDownTime()
	return lastUpAndDownTimeResponseDtoMap
}

func (s *Service) appendStatusTableDataInMap() {
	statusRepository := repository.StatusRepository{}
	statusTableData, _ := statusRepository.GetDataFromDb()
	for _, statusData := range statusTableData {
		statusWithTimeDto := dto.StatusWithTimeDto{
			TimestampUtc: statusData.TimeStampUtc,
			Status:       statusData.Status,
		}
		if _, ok := storeIdToAllTableDataMap[statusData.StoreId]; !ok {
			storeIdToAllTableDataMap[statusData.StoreId] = dto.AllTableDataDto{}
		}
		if len(storeIdToAllTableDataMap[statusData.StoreId].StatusWithTimeList) == 0 {
			storeIdToAllTableDataMap[statusData.StoreId] = dto.AllTableDataDto{StatusWithTimeList: []dto.StatusWithTimeDto{statusWithTimeDto}}
		} else {
			statusWithTimeDtoList := storeIdToAllTableDataMap[statusData.StoreId].StatusWithTimeList
			statusWithTimeDtoList = append(statusWithTimeDtoList, statusWithTimeDto)
			storeIdAllCollectionData := storeIdToAllTableDataMap[statusData.StoreId]
			storeIdAllCollectionData.StatusWithTimeList = statusWithTimeDtoList
			storeIdToAllTableDataMap[statusData.StoreId] = storeIdAllCollectionData
		}
	}
}

func (s *Service) appendBusinessHoursTableDataInMap() {
	businessHoursRepository := repository.BusinessHoursRepository{}
	businessHoursTableData, _ := businessHoursRepository.GetDataFromDb()
	for _, businessHoursData := range businessHoursTableData {
		weekDaysWithStartAndEndTimeDto := dto.WeekDaysWithStartAndEndTime{
			DayOfWeek: businessHoursData.DaysOfWeek,
			StartTime: businessHoursData.StartTimeLocal,
			EndTime:   businessHoursData.EndTimeLocal,
		}
		if _, ok := storeIdToAllTableDataMap[businessHoursData.StoreId]; !ok {
			storeIdToAllTableDataMap[businessHoursData.StoreId] = dto.AllTableDataDto{}
		}
		WeekDaysWithStartAndEndTimeDtoList := storeIdToAllTableDataMap[businessHoursData.StoreId].WeekDaysWithStartAndEndTimeList
		WeekDaysWithStartAndEndTimeDtoList = append(WeekDaysWithStartAndEndTimeDtoList, weekDaysWithStartAndEndTimeDto)
		storeIdAllCollectionData := storeIdToAllTableDataMap[businessHoursData.StoreId]
		storeIdAllCollectionData.WeekDaysWithStartAndEndTimeList = WeekDaysWithStartAndEndTimeDtoList
		storeIdToAllTableDataMap[businessHoursData.StoreId] = storeIdAllCollectionData

	}
}

func (s *Service) appendTimezoneTableDataInMap() {
	timezoneRepository := repository.TimezoneRepository{}
	timezoneTableData, _ := timezoneRepository.GetDataFromDb()
	for _, timezoneData := range timezoneTableData {
		if _, ok := storeIdToAllTableDataMap[timezoneData.StoreId]; !ok {
			storeIdToAllTableDataMap[timezoneData.StoreId] = dto.AllTableDataDto{
				Timezone: timezoneData.TimezoneStr,
			}
		} else {
			storeIdAllCollectionData := storeIdToAllTableDataMap[timezoneData.StoreId]
			storeIdAllCollectionData.Timezone = timezoneData.TimezoneStr
			storeIdToAllTableDataMap[timezoneData.StoreId] = storeIdAllCollectionData
		}
	}
}

func (s *Service) performBusinessLogicOfCalculatingUpAndDownTime() map[string]dto.UpAndDownTimeResponseDto {
	lastUpAndDownTimeResponseDtoMap := make(map[string]dto.UpAndDownTimeResponseDto)
	var currentTime = time.Now().UTC()
	for storeId, data := range storeIdToAllTableDataMap {
		weekDayWithStartTimeMap := make(map[int]time.Time)
		weekDayWithEndTimeMap := make(map[int]time.Time)
		var totalBusinessHoursInWeek int
		var totalBusinessHoursInDay int
		var totalBusinessMinutesInHour = 60
		weekdayWithStatusMap := make(map[int][]dto.StatusWithTimeDto)
		for _, weekDay := range data.WeekDaysWithStartAndEndTimeList {
			weekDayWithStartTimeMap[weekDay.DayOfWeek] = getUtcTime(weekDay.StartTime, data.Timezone)
			weekDayWithEndTimeMap[weekDay.DayOfWeek] = getUtcTime(weekDay.EndTime, data.Timezone)
			currentDayDifference := timeDifference(weekDayWithStartTimeMap[weekDay.DayOfWeek], weekDayWithEndTimeMap[weekDay.DayOfWeek])
			totalBusinessHoursInWeek += int(currentDayDifference.Hours())
			if weekDay.DayOfWeek == constants.ParseIntoDays(currentTime.Weekday().String()) {
				totalBusinessHoursInDay += int(currentDayDifference.Hours())
			}
			fmt.Println(currentTime.Hour(), weekDayWithStartTimeMap[weekDay.DayOfWeek].Hour(), weekDayWithEndTimeMap[weekDay.DayOfWeek].Hour())
			fmt.Println(currentTime, weekDayWithStartTimeMap[weekDay.DayOfWeek], weekDayWithEndTimeMap[weekDay.DayOfWeek])
			if totalBusinessMinutesInHour == 0 && weekDayWithStartTimeMap[weekDay.DayOfWeek].Hour() == currentTime.Hour() && weekDayWithEndTimeMap[weekDay.DayOfWeek].Hour() == currentTime.Hour() {
				totalBusinessMinutesInHour -= weekDayWithStartTimeMap[weekDay.DayOfWeek].Minute()
				totalBusinessMinutesInHour -= 60 - weekDayWithEndTimeMap[weekDay.DayOfWeek].Minute()
			}
		}
		for _, statusAndTimeData := range data.StatusWithTimeList {
			currentDay := int(statusAndTimeData.TimestampUtc.Weekday())
			weekdayWithStatusMap[currentDay-1] = append(weekdayWithStatusMap[currentDay-1], statusAndTimeData)
		}
		var statusWithTimeInActiveBusinessHoursList []dto.StatusWithTimeDto
		for weekdayNumber, weekdayWithStatusList := range weekdayWithStatusMap {
			lastStatusType := ""
			for _, weekdayWithStatus := range weekdayWithStatusList {
				timeStamp := getUtcTimeInFormat(weekdayWithStatus.TimestampUtc)
				fmt.Println(timeStamp, weekDayWithStartTimeMap[weekdayNumber], weekDayWithEndTimeMap[weekdayNumber])
				fmt.Println(timeStamp.Unix(), weekDayWithStartTimeMap[weekdayNumber].Unix(), weekDayWithEndTimeMap[weekdayNumber].Unix())
				if timeStamp.Unix() > weekDayWithEndTimeMap[weekdayNumber].Unix() {
					continue
				} else if timeStamp.Unix() < weekDayWithStartTimeMap[weekdayNumber].Unix() {
					lastStatusType = weekdayWithStatus.Status
					continue
				} else if lastStatusType == constants.InActiveStatusString {
					statusWithTimeInActiveBusinessHoursList = []dto.StatusWithTimeDto{{Status: lastStatusType, TimestampUtc: weekDayWithStartTimeMap[weekdayNumber]}}
					lastStatusType = ""
				}
				statusWithTimeInActiveBusinessHoursList = append(statusWithTimeInActiveBusinessHoursList, weekdayWithStatus)
			}
			statusWithTimeInActiveBusinessHoursList = append(statusWithTimeInActiveBusinessHoursList, dto.StatusWithTimeDto{Status: "active", TimestampUtc: weekDayWithEndTimeMap[weekdayNumber]})
		}
		totalBusinessHoursInWeek += (7 - len(data.WeekDaysWithStartAndEndTimeList)) * 24
		lastDownTimeForWeek := int(s.inactiveDurationForLastWeek(statusWithTimeInActiveBusinessHoursList).Hours())
		lastUpTimeForWeek := totalBusinessHoursInWeek - lastDownTimeForWeek
		lastDownTimeForDay := int(s.inactiveDurationForLastDay(statusWithTimeInActiveBusinessHoursList).Hours())
		lastUpTimeForDay := totalBusinessHoursInDay - lastDownTimeForDay
		lastDownTimeForHour := int(s.inactiveDurationForLastHour(statusWithTimeInActiveBusinessHoursList).Minutes())
		lastUpTimeForHour := totalBusinessMinutesInHour - lastDownTimeForHour
		lastUpAndDownTimeResponse := dto.UpAndDownTimeResponseDto{
			DownTimeLastWeek: strconv.Itoa(lastDownTimeForWeek),
			UpTimeLastWeek:   strconv.Itoa(lastUpTimeForWeek),
			DownTimeLastDay:  strconv.Itoa(lastDownTimeForDay),
			UpTimeLastDay:    strconv.Itoa(lastUpTimeForDay),
			DownTimeLastHour: strconv.Itoa(lastDownTimeForHour),
			UpTimeLastHour:   strconv.Itoa(lastUpTimeForHour),
		}
		lastUpAndDownTimeResponseDtoMap[storeId] = lastUpAndDownTimeResponse
	}
	return lastUpAndDownTimeResponseDtoMap
}

func (s *Service) inactiveDurationForLastWeek(statusWithTimeDtoList []dto.StatusWithTimeDto) time.Duration {
	var inactiveDuration time.Duration
	for index := 0; index+1 < len(statusWithTimeDtoList); index++ {
		statusWithTime := statusWithTimeDtoList[index]
		nextStatusWithTime := statusWithTimeDtoList[index+1]
		status := statusWithTime.Status
		timestamp := statusWithTime.TimestampUtc
		if status == constants.ActiveStatusString {
			continue
		}
		timeDiff := timeDifference(timestamp, nextStatusWithTime.TimestampUtc)
		inactiveDuration += timeDiff
	}
	return inactiveDuration
}

func (s *Service) inactiveDurationForLastDay(statusWithTimeDtoList []dto.StatusWithTimeDto) time.Duration {
	curTime := time.Now().UTC()
	var inactiveDuration time.Duration
	for index := 0; index+1 < len(statusWithTimeDtoList); index++ {
		statusWithTime := statusWithTimeDtoList[index]
		nextStatusWithTime := statusWithTimeDtoList[index+1]
		status := statusWithTime.Status
		timestamp := statusWithTime.TimestampUtc
		if status == constants.ActiveStatusString || timestamp.Before(curTime.Add(-24*time.Hour)) {
			continue
		}
		timeDiff := timeDifference(timestamp, minTimestamp(nextStatusWithTime.TimestampUtc, curTime))
		inactiveDuration += timeDiff
	}
	return inactiveDuration
}

func (s *Service) inactiveDurationForLastHour(statusWithTimeDtoList []dto.StatusWithTimeDto) time.Duration {
	curTime := time.Now().UTC()
	var inactiveDuration time.Duration
	for index := 0; index+1 < len(statusWithTimeDtoList); index++ {
		statusWithTime := statusWithTimeDtoList[index]
		nextStatusWithTime := statusWithTimeDtoList[index+1]
		status := statusWithTime.Status
		timestamp := statusWithTime.TimestampUtc
		if status == constants.ActiveStatusString || timestamp.Before(curTime.Add(-time.Hour)) {
			continue
		}
		timeDiff := timeDifference(timestamp, minTimestamp(nextStatusWithTime.TimestampUtc, curTime))
		inactiveDuration += timeDiff
	}
	return inactiveDuration
}

func getUtcTimeInFormat(currentTime time.Time) time.Time {
	currentTimeInString := currentTime.String()
	return getUtcTime(currentTimeInString[11:19], "UTC")
}

func getUtcTime(currentTime string, currentTimezone string) time.Time {
	location, err := time.LoadLocation(currentTimezone)
	if err != nil {
		fmt.Println(err)
		return time.Time{}
	}
	t, err := time.ParseInLocation(constants.HoursMinutesSecondsFormat, currentTime, location)
	if err != nil {
		fmt.Println(err)
		return time.Time{}
	}
	utcLocation, err := time.LoadLocation("UTC")
	return t.In(utcLocation)
}

func timeDifference(startTime time.Time, endTime time.Time) time.Duration {

	diff := maxTimestamp(startTime, endTime).Sub(minTimestamp(startTime, endTime))
	return diff
}

func minTimestamp(time1 time.Time, time2 time.Time) time.Time {
	if time1.Before(time2) {
		return time1
	}
	return time2
}

func maxTimestamp(time1 time.Time, time2 time.Time) time.Time {
	if time1.Before(time2) {
		return time2
	}
	return time1
}
