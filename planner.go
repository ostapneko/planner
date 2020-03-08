package planner

import (
	"fmt"
)

type developerId string

type durationDays int

type day int

type task struct {
	name         string
	attributions map[developerId]durationDays
}

type developer struct {
	id      developerId
	offDays []day
}

type calendar struct {
	days []day
}

type supportWeek struct {
	firstDay day
	lastDay  day
	devId developerId
}

// tasks are sorted in priority order: highest priority first
func checkGraph(tasks []task, developers []developer, supportWeeks []supportWeek, cal calendar) error {
	devMap := make(map[developerId]developer, len(developers))
	for _, dev := range developers {
		devMap[dev.id] = dev
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
func checkSupportWeeks(supportWeeks []supportWeek, devMap map[developerId]developer) error {
	minWeek := day(1e6)
	maxWeek := day(0)
	for _, week := range supportWeeks {
		if week.firstDay < minWeek {
			minWeek = week.firstDay
		}

		if week.lastDay > maxWeek {
			maxWeek = week.lastDay
		}

		if _, prs := devMap[week.devId]; !prs {
			return fmt.Errorf("developer %s mentioned in support week %v does not exit", week.devId, week)
		}
	}

	allDays := make([]bool, maxWeek-minWeek+1)
	for _, week := range supportWeeks {
		isEmpty := true
		for i := week.firstDay; i < week.lastDay+1; i++ {
			isEmpty = false
			if allDays[i-minWeek] {
				return fmt.Errorf("day %d is in more than one week", i)
			}
			allDays[i-minWeek] = true
		}
		if isEmpty {
			return fmt.Errorf("support week %v is empty", week)
		}
	}
	return nil
}

func checkDevAttributions(tasks []task, devMap map[developerId]developer) error {
	for _, t := range tasks {
		for devId, _ := range t.attributions {
			if _, prs := devMap[devId]; !prs {
				return fmt.Errorf("developer %s mentioned in task %v does not exist", devId, t)
			}
		}
	}
	return nil
}
