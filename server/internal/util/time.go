package util

import "time"

func TodayYYYYMMDD() string {
	return time.Now().Format("20060102")
}
