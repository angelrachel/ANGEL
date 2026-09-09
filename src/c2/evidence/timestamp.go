package evidence

import "time"

type TimestampResult struct {
	UTC  string
	Unix int64
}

func GetTimestamp() TimestampResult {
	now := time.Now().UTC()
	return TimestampResult{UTC: now.Format(time.RFC3339), Unix: now.Unix()}
}
