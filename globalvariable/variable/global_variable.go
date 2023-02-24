package variable

import (
	"time"
)

// Get Current Time
var currentTime = time.Now()

// Format the time
var FormatedTimeiso8601 = currentTime.Format(time.RFC3339)

var expiry = currentTime.Add(time.Hour * 24 * 3)
var thirtyDaysAgo = time.Now().Add(-30 * 24 * time.Hour)


var ExpiryStrFormatted = expiry.Format("2006-01-02 15:04")

var FormattedTimeNowYYYYMMDDHHMM = currentTime.Format("2006-01-02 15:04")

var FormatedTime30DayAgoIso8601 = thirtyDaysAgo.Format(time.RFC3339)
