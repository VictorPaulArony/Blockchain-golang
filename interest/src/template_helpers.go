package helpers

import "time"

// unixToTime converts a Unix timestamp to a formatted string.
func UnixToTime(unixTime int64) string {
	return time.Unix(unixTime, 0).Format("2006-01-02 15:04:05")
}

func Mul(a, b float64) float64 { // New multiplication function
	return a * b
}
