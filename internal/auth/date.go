package auth

import (
	"net/http"
	"strings"
	"time"
)

const ISO8601 = "20060102T150405Z0700"

var (
	dateHeaders    = []string{strings.ToLower(AmzDate), AmzDate, hdrDate}
	errInvalidDate = MissingSecurityHeader("AWS authentication requires a valid Date or x-amz-date header")
)

func getDate(values Values) (date time.Time, err error) {
	if value := getDateString(values); value != "" {
		return parseDate(value)
	}
	return time.Time{}, errInvalidDate
}

func getDateString(values Values) (date string) {
	for _, header := range dateHeaders {
		if date = values.Get(header); date != "" {
			break
		}
	}
	return
}

func parseDate(value string) (date time.Time, err error) {
	date, err = http.ParseTime(value)
	if err != nil {
		date, err = time.Parse(ISO8601, value)
	}
	if err != nil {
		err = errInvalidDate
	}
	return
}
