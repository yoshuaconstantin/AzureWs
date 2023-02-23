package globalvariable

import (

	"time"
)

// Get Current Time
var currentTime = time.Now()

// Format the time
var FormatedTime = currentTime.Format(time.RFC3339)

var expiry = currentTime.Add(time.Hour * 24 * 3)

var ExpiryStr = expiry.Format("2006-01-02 15:04")

