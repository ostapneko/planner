package planner

import (
	"fmt"
)

type DeveloperId string

type DurationDays int

type Day int

type Task struct {
	Name         string
	Attributions map[DeveloperId]DurationDays
}

type Developer struct {
	Id      DeveloperId
	OffDays []Day
}

type Calendar struct {
	Days []Day
}

type SupportWeek struct {
	FirstDay Day
	LastDay  Day
	DevId    DeveloperId
}

// tasks are sorted in priority order: highest priority first
func CheckGraph(tasks []Task, developers []Developer, supportWeeks []SupportWeek, cal Calendar) error {
	devMap := make(map[DeveloperId]Developer, len(developers))
	for _, dev := range developers {
		devMap[dev.Id] = dev
	}

	err := checkDevAttributions(tasks, devMap)
	if err != nil {
		return err
	}

	err = checkSupportWeeks(supportWeeks, devMap)
	if err != nil {
		return err
	}

	return nil
}

//check devs in support weeks exist
// check support weeks are not overlapping, and that week are not empty
func checkSupportWeeks(supportWeeks []SupportWeek, devMap map[DeveloperId]Developer) error {
	minWeek := Day(1e6)
	maxWeek := Day(0)
	for _, week := range supportWeeks {
		if week.FirstDay < minWeek {
			minWeek = week.FirstDay
		}

		if week.LastDay > maxWeek {
			maxWeek = week.LastDay
		}

		if _, prs := devMap[week.DevId]; !prs {
			return fmt.Errorf("Developer %s mentioned in support week %v does not exit", week.DevId, week)
		}
	}

	allDays := make([]bool, maxWeek-minWeek+1)
	for _, week := range supportWeeks {
		isEmpty := true
		for i := week.FirstDay; i < week.LastDay+1; i++ {
			isEmpty = false
			if allDays[i-minWeek] {
				return fmt.Errorf("Day %d is in more than one week", i)
			}
			allDays[i-minWeek] = true
		}
		if isEmpty {
			return fmt.Errorf("support week %v is empty", week)
		}
	}
	return nil
}

func checkDevAttributions(tasks []Task, devMap map[DeveloperId]Developer) error {
	for _, t := range tasks {
		for devId, _ := range t.Attributions {
			if _, prs := devMap[devId]; !prs {
				return fmt.Errorf("Developer %s mentioned in Task %v does not exist", devId, t)
			}
		}
	}
	return nil
}
