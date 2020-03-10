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

	err := checkTasks(planning.Tasks, devMap)
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

func checkTasks(tasks []*Task, devMap map[DeveloperId]*Developer) error {
	for _, t := range tasks {
		var latestLastDay Day = 0
		for devId, attr := range t.Attributions {
			if _, prs := devMap[devId]; !prs {
				return fmt.Errorf("developer %s mentioned in Task %v does not exist", devId, t)
			}

			if attr.LastDay != nil && attr.FirstDay != nil && int(attr.DurationDays) != int(*attr.LastDay-*attr.FirstDay+1) {
				return fmt.Errorf("duration is inconsistent for attribution %+v of task %s", *attr, t.Name)
			}

			if attr.LastDay != nil && attr.FirstDay == nil {
				return fmt.Errorf("attribution %+v of task %s has a last day but no first day", *attr, t.Name)
			}

			if attr.LastDay != nil && *attr.LastDay > latestLastDay {
				latestLastDay = *attr.LastDay
			}
		}
		if latestLastDay == 0 && t.LastDay != nil {
			return fmt.Errorf("task has a last day but no attribution has a last day")
		}
	}
	return nil
}
