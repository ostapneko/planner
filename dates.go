package planner

import (
	"fmt"
	"math"
	"time"
)

// from 15/02/2020 -> 18307 (nb of days since epoch)
func DateToDay(str string) (Day, error) {
	t, err := time.Parse("02/01/2006", str)

	if err != nil {
		return 0, fmt.Errorf("error parsing date %s, should be in the format 25/05/1983", str)
	}

	fromEpoch := t.Sub(time.Unix(0, 0))

	return Day(math.Floor(fromEpoch.Hours() / 24.0)), nil
}
