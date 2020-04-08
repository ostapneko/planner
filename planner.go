package planner

import (
	"fmt"
	"time"
)

type Planning struct {
	StartDay     Day
	Holidays     Days
	Developers   []*Developer
	SupportWeeks []*SupportWeek `yaml:"supportWeeks"`
	// tasks are sorted in priority order: highest priority first
	Tasks []*Task `yaml:"tasks"`
}

type DeveloperId string

type EffortDays int

type Day int
type Days []Day

func (d Days) Len() int {
	return len(d)
}

func (d Days) Less(i, j int) bool {
	return d[i] < d[j]
}

func (d Days) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

type Task struct {
	Name         string
	Attributions map[DeveloperId]*Attribution
	LastDay      *Day
}

type Attribution struct {
	EffortDays EffortDays
	FirstDay   *Day
	LastDay    *Day
}

type Developer struct {
	Id      DeveloperId
	OffDays Days `yaml:"offDays"`
}

type SupportWeek struct {
	FirstDay Day
	LastDay  Day
	DevId    DeveloperId `yaml:"devId"`
}

func CheckPlanning(planning *Planning) error {
	devMap := make(map[DeveloperId]*Developer, len(planning.Developers))
	for _, dev := range planning.Developers {
		devMap[dev.Id] = dev
	}

	holidaysMap := make(map[Day]interface{}, len(planning.Holidays))
	for _, day := range planning.Holidays {
		holidaysMap[day] = nil
	}

	err := checkTasks(planning.Tasks, devMap, holidaysMap, planning.SupportWeeks)
	if err != nil {
		return err
	}

	err = checkSupportWeeks(planning.SupportWeeks, devMap)
	if err != nil {
		return err
	}

	return nil
}

// ForecastCompletion attributes a FirstDay and a LastDay to all attributions,
// as well as a last day to all tasks
func ForecastCompletion(planning *Planning) {
	// This maps developers to all their non-worked days
	devToOffDays := make(map[DeveloperId]map[Day]bool)

	// fill the map with empty maps
	for _, developer := range planning.Developers {
		devToOffDays[developer.Id] = make(map[Day]bool)
	}

	// holidays
	for _, holiday := range planning.Holidays {
		for _, developer := range planning.Developers {
			devToOffDays[developer.Id][holiday] = true
		}
	}

	// off days
	for _, developer := range planning.Developers {
		for _, day := range developer.OffDays {
			devToOffDays[developer.Id][day] = true
		}
	}

	// support weeks
	for _, week := range planning.SupportWeeks {
		for i := week.FirstDay; i <= week.LastDay; i++ {
			devToOffDays[week.DevId][i] = true
		}
	}


	// devToLatestDay associate a the latest day that was allocated for each developer
	// as we go through each task and each attribution by order of priority, we are going to increment this day
	// until we find a non-holiday, non-off-day, non-support-week-day, non-week ends for this developer, and repeat until
	// all the effort days for all attributions have been fullfilled
	devToLatestDay := make(map[DeveloperId]Day)

	for _, developer := range planning.Developers {
		devToLatestDay[developer.Id] = planning.StartDay
	}

	for _, task := range planning.Tasks {
		var lastTaskDay *Day
		task.LastDay = nil
		for developerId, attribution := range task.Attributions {
			attribution.FirstDay = nil
			attribution.LastDay = nil
			var effort EffortDays = 0
			for effort < attribution.EffortDays {
				day := devToLatestDay[developerId]
				// if the day is not off, increment the effort
				if _, prs := devToOffDays[developerId][day]; !prs && !isWeekEnd(day) {
					effort++
					// if the first day is not set, set it
					if attribution.FirstDay == nil {
						firstDay := devToLatestDay[developerId]
						attribution.FirstDay = &firstDay
					}
				}

				devToLatestDay[developerId] = day + 1
			}
			attrLastDay := devToLatestDay[developerId] - 1
			attribution.LastDay = &attrLastDay

			if lastTaskDay == nil || attrLastDay > *lastTaskDay {
				lastTaskDay = &attrLastDay
			}
		}
		task.LastDay = lastTaskDay
	}
}

func isWeekEnd(day Day) bool {
	weekDay := DayToTime(day).Weekday()
	return weekDay == time.Saturday || weekDay == time.Sunday
}

// check devs in support weeks exist
// check support weeks are not overlapping, and that weeks are not empty
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

func checkTasks(tasks []*Task, devMap map[DeveloperId]*Developer, calendarMap map[Day]interface{}, supportWeeks []*SupportWeek) error {
	for _, t := range tasks {
		var latestLastDay Day = 0
		var allAttributionsHaveLastDay = true
		for devId, attr := range t.Attributions {
			dev, devPrs := devMap[devId]
			if !devPrs {
				return fmt.Errorf("developer %s mentioned in Task %v does not exist", devId, t)
			}

			devSupportWeeks := make([]*SupportWeek, 0)
			for _, w := range supportWeeks {
				if w.DevId == devId {
					devSupportWeeks = append(devSupportWeeks, w)
				}
			}

			if attr.LastDay != nil && attr.FirstDay != nil {
				computedEffortDays := effort(*attr.FirstDay, *attr.LastDay, calendarMap, dev.OffDays, devSupportWeeks)
				if attr.EffortDays != computedEffortDays {
					return fmt.Errorf("effort is inconsistent for attribution %+v of task %s. Should be %d, but got %d", *attr, t.Name, computedEffortDays, attr.EffortDays)
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

func effort(firstDay Day, lastDay Day, holidays map[Day]interface{}, offDays []Day, weeks []*SupportWeek) EffortDays {
	res := 0
	for i := firstDay; i < lastDay+1; i++ {
		// skip holidays not in calendar
		if _, prs := holidays[i]; prs {
			continue
		}

		// skip off days
		isOffDay := false
		for _, d := range offDays {
			if d == i {
				isOffDay = true
			}
		}

		if isOffDay {
			continue
		}

		// skip support days
		isSupportDay := false
		for _, w := range weeks {
			if i >= w.FirstDay && i <= w.LastDay {
				isSupportDay = true
			}
		}

		if isSupportDay {
			continue
		}

		res++
	}

	return EffortDays(res)
}
