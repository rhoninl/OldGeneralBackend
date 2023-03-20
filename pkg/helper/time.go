package helper

import "time"

func GetTimeDiff(timestamp int64) int64 {
	return int64(time.UnixMicro(timestamp).Sub(time.Now()).Hours() / 24)
}
