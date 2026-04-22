package task

import (
	"time"

	"github.com/teambition/rrule-go"
)

func GenerateOccurrences(rruleStr string, from, to time.Time) ([]time.Time, error) {
	if rruleStr == "" {
		return nil, nil
	}

	fullRRule := "DTSTART:" + from.Format("20060102T000000Z") + "\n" + rruleStr

	rule, err := rrule.StrToRRule(fullRRule)
	if err != nil {
		return nil, err
	}

	fromUTC := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)
	toUTC := time.Date(to.Year(), to.Month(), to.Day(), 23, 59, 59, 0, time.UTC)

	return rule.Between(fromUTC, toUTC, true), nil
}
