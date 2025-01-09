package utils

import "time"

func NowString() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}
