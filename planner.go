package planner

import (
	"fmt"
)

type Planning struct {
	Calendar     []Day
	Developers   []*Developer
	SupportWeeks []*SupportWeek `yaml:"supportWeeks"`
	Tasks        []*Task        `yaml:"tasks"`
}

type DeveloperId string

type DurationDays int

type Day int

type Task struct {
	Name         string
	Attributions map[DeveloperId]*Attribution
	LastDay      *Day
}

type Attribution struct {
	DurationDays DurationDays
	FirstDay     *Day
	LastDay      *Day
}

type Developer struct {
	Id      DeveloperId
	OffDays []Day `yaml:"offDays"`
}

type SupportWeek struct {
	FirstDay Day
	LastDay  Day
	DevId    DeveloperId `yaml:"devId"`
}

// tasks are sorted in priority order: highest priority first
func CheckPlanning(planning *Planning) error {
	devMap := make(map[DeveloperId]*Developer, len(planning.Developers))
	for _, dev := range planning.Developers {
		devMap[dev.Id] = dev
	}

	calendarMap := make(map[Day]interface{}, len(planning.Calendar))
	for _, day := range planning.Calendar {
		calendarMap[day] = nil
	}

	err := checkTasks(planning.Tasks, devMap, calendarMap)
	if err != nil {
		return err
	}

	err = checkSupportWeeks(planning.SupportWeeks, devMap)
	if err != nil {
		return err
	}

	return nil
}

//check devs in support weeks exist
// check support weeks are not overlapping, and that week are not empty
func checkSupportWeeks(supportWeeks []*SupportWeek, devMap map[DeveloperId]*Developer) error {
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
			return fmt.Errorf("developer %s mentioned in support week %v does not exit", week.DevId, week)
		}
	}

	allDays := make([]bool, maxWeek-minWeek+1)
	for _, week := range supportWeeks {
		isEmpty := true
		for i := week.FirstDay; i < week.LastDay+1; i++ {
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

func checkTasks(tasks []*Task, devMap map[DeveloperId]*Developer, calendarMap map[Day]interface{}) error {
	for _, t := range tasks {
		var latestLastDay Day = 0
		var allAttributionsHaveLastDay = true
		for devId, attr := range t.Attributions {
			dev, devPrs := devMap[devId]
			if !devPrs {
				return fmt.Errorf("developer %s mentioned in Task %v does not exist", devId, t)
			}

			if attr.LastDay != nil && attr.FirstDay != nil {
				computedDurationDays := duration(*attr.FirstDay, *attr.LastDay, calendarMap, dev.OffDays)
				if attr.DurationDays != computedDurationDays {
					return fmt.Errorf("duration is inconsistent for attribution %+v of task %s. Should be %d, but got %d", *attr, t.Name, computedDurationDays, attr.DurationDays)
				}
			}

			if attr.LastDay != nil && attr.FirstDay == nil {
				return fmt.Errorf("attribution %+v of task %s has a last day but no first day", *attr, t.Name)
			}

			if attr.LastDay != nil && *attr.LastDay > latestLastDay {
				latestLastDay = *attr.LastDay
			}

			if attr.LastDay == nil {
				allAttributionsHaveLastDay = false
			}
		}

		if latestLastDay == 0 && t.LastDay != nil {
			return fmt.Errorf("task %s has a last day but no attribution has a last day", t.Name)
		}

		if allAttributionsHaveLastDay && t.LastDay != nil && *t.LastDay != latestLastDay {
			return fmt.Errorf("task %s has a LastDay inconsistent with its attributions", t.Name)
		}
	}
	return nil
}

func duration(firstDay Day, lastDay Day, calendar map[Day]interface{}, offDays []Day) DurationDays {
	res := 0
	for i := firstDay; i < lastDay+1; i++ {
		// skip days not in calendar
		if _, prs := calendar[i]; !prs {
			continue
		}

		// skip off days
		for _, d := range offDays {
			isOffDay := false
			if d == i {
				isOffDay = true
			}

			if isOffDay {
				continue
			}
		}

		res++
	}

	return DurationDays(res)
}
