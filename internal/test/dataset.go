package test

import (
	"fmt"
	"time"

	"fetchrewards.com/points-api/internal/model"
)

var Data = []model.Transaction{
	{
		Payer:     "DANNON",
		Points:    1000,
		Timestamp: ParseTime("2020-11-02T14:00:00Z"),
	},
	{
		Payer:     "UNILEVER",
		Points:    200,
		Timestamp: ParseTime("2020-10-31T11:00:00Z"),
	},
	{
		Payer:     "DANNON",
		Points:    -200,
		Timestamp: ParseTime("2020-10-31T15:00:00Z"),
	},
	{
		Payer:     "MILLER COORS",
		Points:    10000,
		Timestamp: ParseTime("2020-11-01T14:00:00Z"),
	},
	{
		Payer:     "DANNON",
		Points:    300,
		Timestamp: ParseTime("2020-10-31T10:00:00Z"),
	},
}

func ParseTime(timestampStr string) time.Time {
	result, err := time.Parse(time.RFC3339, timestampStr)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse time %s", timestampStr))
	}
	return result
}
