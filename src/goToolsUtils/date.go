package goToolsUtils

import "time"

func GetCurrentTime() string {
	return Date(time.Now().UnixNano() / 1e6)
}

func Date(million int64) string {
	return time.Unix(million/1e3, million%1e3).Format("2006-01-02 15:04:05")
}
