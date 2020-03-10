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
