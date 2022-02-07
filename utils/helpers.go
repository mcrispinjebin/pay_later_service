package utils

import "time"

func CurrentMillis() int {
	return int(time.Now().UnixNano() / 1e6)
}
