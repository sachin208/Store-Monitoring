package constants

const (
	DbConnectionString        = "root:Nonu5653@@1@tcp(localhost:3306)/Store_Monitoring?parseTime=true"
	StatusFilePath            = "/Users/sachin/Desktop/zest/go/src/store_monitoring/db/file/Store_Status.csv"
	TimezoneFilePath          = "/Users/sachin/Desktop/zest/go/src/store_monitoring/db/file/Timezone.csv"
	BusinessHoursFilePath     = "/Users/sachin/Desktop/zest/go/src/store_monitoring/db/file/Menu_Hours.csv"
	FileFolderPath            = "/Users/sachin/Desktop/zest/go/src/store_monitoring/db/file/reports/"
	ActiveStatusString        = "active"
	InActiveStatusString      = "inactive"
	HoursMinutesSecondsFormat = "15:04:05"
)

func ParseIntoDays(day string) int {
	switch day {
	case "Monday":
		return 0
	case "Tuesday":
		return 1
	case "Wednesday":
		return 2
	case "Thursday":
		return 3
	case "Friday":
		return 4
	case "Saturday":
		return 5
	case "Sunday":
		return 6
	}
	return -1
}
