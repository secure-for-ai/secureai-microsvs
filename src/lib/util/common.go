package util

import "time"

func GetNowTimestamp() int64 {
	return time.Now().Unix()
}
