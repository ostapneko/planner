package planner

import "testing"

func TestDateToDay(t *testing.T) {
	act, err := DateToDay("02/01/1970")

	if err != nil {
		t.Error(err)
	}

	if act != 1 {
		t.Errorf("exp 1, got %d", act)
	}
}

func TestDayTodate(t *testing.T) {
	act := DayToDate(18307)

	exp := "15/02/2020"

	if act != exp {
		t.Errorf("exp %s, got %s", exp, act)
	}
}
