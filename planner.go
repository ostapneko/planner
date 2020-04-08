package planner

import (
	"fmt"
	"sort"
)

type Planning struct {
	Calendar     Days
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

	calendarMap := make(map[Day]interface{}, len(planning.Calendar))
	for _, day := range planning.Calendar {
		calendarMap[day] = nil
	}

	err := checkTasks(planning.Tasks, devMap, calendarMap, planning.SupportWeeks)
	if err != nil {
		return err
	}

	err = checkSupportWeeks(planning.SupportWeeks, devMap)
	if err != nil {
		return err
	}

	return nil
}

func ForecastCompletion(planning *Planning, now Day) {
	devToDays := make(map[DeveloperId][]Day)
	for _, developer := range planning.Developers {
		devToDays[developer.Id] = availableDays(developer, planning.Calendar, planning.SupportWeeks, now)
	}

	for _, task := range planning.Tasks {
		// skip finished tasks
		if task.LastDay != nil && *task.LastDay < now {
			continue
		}

		// for each attribution, add days from allocatable developer days
		// until those days are equal to the attribution's effort day, or until we run out
		// along the way, we bump the last attribution day and assign it to the attributions in case there
		// is enough allocated days to fulfill the attribution
		for developerId, attribution := range task.Attributions {
			var attributedEffort EffortDays = 0

			var lastAttrDay *Day
			if len(devToDays[developerId]) > 0 {
				attribution.FirstDay = &devToDays[developerId][0]
			}

			for len(devToDays[developerId]) > 0 && attributedEffort < attribution.EffortDays {
				// pop the earliest day
				lastAttrDay = &devToDays[developerId][0]
				devToDays[developerId] = devToDays[developerId][1:]
				attributedEffort += 1
			}

			if attributedEffort == attribution.EffortDays {
				attribution.LastDay = lastAttrDay
			}

		}

		// if all attributions are fulfilled, update the task's last day
		var lastDay *Day
		allAttrAreFulfilled := true
		for _, attribution := range task.Attributions {
			if attribution.LastDay == nil {
				allAttrAreFulfilled = false
				break
			}

			if lastDay == nil || *attribution.LastDay > *lastDay {
				lastDay = attribution.LastDay
			}
		}

		if allAttrAreFulfilled {
			task.LastDay = lastDay
		}
	}
}

func availableDays(dev *Developer, cal Days, weeks []*SupportWeek, now Day) []Day {
	res := make(Days, 0)
	offDays := map[Day]bool{}
	for _, day := range dev.OffDays {
		offDays[day] = true
	}
	supportDays := supportWeekDays(weeks)

	for _, day := range cal {
		_, isOff := offDays[day]
		_, isSupport := supportDays[day]
		isInFuture := day >= now
		if !isOff && !isSupport && isInFuture {
			res = append(res, day)
		}
	}
	sort.Sort(res)
	return res
}

func supportWeekDays(weeks []*SupportWeek) map[Day]bool {
	res := map[Day]bool{}
	for _, week := range weeks {
		for i := week.FirstDay; i <= week.LastDay; i++ {
			res[i] = true
		}
	}
	return res
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

func effort(firstDay Day, lastDay Day, calendar map[Day]interface{}, offDays []Day, weeks []*SupportWeek) EffortDays {
	res := 0
	for i := firstDay; i < lastDay+1; i++ {
		// skip days not in calendar
		if _, prs := calendar[i]; !prs {
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
