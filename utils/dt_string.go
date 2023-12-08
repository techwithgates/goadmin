package utils

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func dateTimeToString(dataType string, value any) string {
	if dataType == "date" {
		// convert time.Time to string date format to be compatible with <input type='date'>
		if date, ok := value.(time.Time); ok {
			value = date.Format("2006-01-02")
		}
	}

	if dataType == "timestamp without time zone" {
		// convert time.Time to string datetime format to be compatible with <input type='datetime-local'>
		if dateTime, ok := value.(time.Time); ok {
			value = dateTime.Format("2006-01-02T15:04")
		}
	}

	if dataType == "time without time zone" {
		// convert pgtype.Time to string time format to be compatible with <input type='time'>
		if time, ok := value.(pgtype.Time); ok {
			time, _ := time.Value()
			if time, ok := time.(string); ok {
				value = time[:5]
			}
		}
	}

	if result, ok := value.(string); ok {
		return result
	}

	return ""
}
