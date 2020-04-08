package planner

import (
	"fmt"
	"math"
	"time"
)

var dateFormat = "02/01/2006"

// from 15/02/2020 -> 18307 (nb of days since epoch)
func DateToDay(str string) (Day, error) {
	t, err := time.Parse(dateFormat, str)

	if err != nil {
		return 0, fmt.Errorf("error parsing date %s, should be in the format 25/05/1983", str)
	}

	fromEpoch := t.Sub(time.Unix(0, 0))

	return Day(math.Floor(fromEpoch.Hours() / 24.0)), nil
}

// from  18307 (nb of days since epoch) -> 15/02/2020
func DayToDate(day Day) string {
	t := DayToTime(day)
	return t.Format(dateFormat)
}

func DayToTime(day Day) time.Time {
	epoch := time.Unix(0, 0)
	return epoch.Add(time.Duration(int(day)) * time.Hour * 24)
}
