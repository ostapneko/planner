package planner

import "fmt"

type PlanningInput struct {
	StartDay     string `yaml:"startDay"`
	Holidays     []string
	Developers   []*DeveloperInput   `yaml:"developers"`
	SupportWeeks []*SupportWeekInput `yaml:"supportWeeks"`
	Tasks        []*TaskInput        `yaml:"tasks"`
}

type TaskInput struct {
	Name         string
	Attributions map[DeveloperId]*AttributionInput
}

type AttributionInput struct {
	Effort   EffortDays
	FirstDay *string `yaml:"firstDay"`
	LastDay  *string `yaml:"lastDay"`
}

type SupportWeekInput struct {
	FirstDay string      `yaml:"firstDay"`
	LastDay  string      `yaml:"lastDay"`
	DevId    DeveloperId `yaml:"devId"`
}

type DeveloperInput struct {
	Id      DeveloperId
	OffDays []string `yaml:"offDays"`
	Starts  *string  `yaml:"starts,omitempty"`
	Leaves    *string  `yaml:"leaves,omitempty"`
	Utilization float64 `yaml:"utilization"`
}

func NewPlanning(input PlanningInput) (*Planning, error) {
	holidays := make([]Day, len(input.Holidays))
	for i, s := range input.Holidays {
		d, err := DateToDay(s)
		if err != nil {
			return nil, fmt.Errorf("error parsing cal %s", err)
		}
		holidays[i] = d
	}

	devs := make([]*Developer, len(input.Developers))

	for i, input := range input.Developers {
		dev, err := newDeveloper(input)
		if err != nil {
			return nil, fmt.Errorf("error parsing developer %s", err)
		}
		devs[i] = dev
	}

	weeks := make([]*SupportWeek, len(input.SupportWeeks))
	for i, input := range input.SupportWeeks {
		week, err := newSupportWeek(input)
		if err != nil {
			return nil, fmt.Errorf("error parsing week %s", err)
		}
		weeks[i] = week
	}

	tasks := make([]*Task, len(input.Tasks))
	for i, input := range input.Tasks {
		task, err := newTask(input)
		if err != nil {
			return nil, fmt.Errorf("error parsing task %s", err)
		}
		tasks[i] = task
	}

	startDay, err := DateToDay(input.StartDay)
	if err != nil {
		return nil, fmt.Errorf("error parsing start day: %s", err)
	}

	return &Planning{
		StartDay:     startDay,
		Holidays:     holidays,
		Developers:   devs,
		SupportWeeks: weeks,
		Tasks:        tasks,
	}, nil
}

func newDeveloper(input *DeveloperInput) (*Developer, error) {
	offDays := make([]Day, len(input.OffDays))

	for i, s := range input.OffDays {
		d, err := DateToDay(s)
		if err != nil {
			return nil, err
		}
		offDays[i] = d
	}

	var starts *Day
	if input.Starts != nil {
		day, err := DateToDay(*input.Starts)
		if err != nil {
			return nil, err
		}
		starts = &day
	}

	var leaves *Day
	if input.Leaves != nil {
		day, err := DateToDay(*input.Leaves)
		if err != nil {
			return nil, err
		}
		leaves = &day
	}

	return &Developer{
		Id:      input.Id,
		OffDays: offDays,
		Starts:  starts,
		Leaves:    leaves,
		Utilization: input.Utilization,
	}, nil
}

func newSupportWeek(input *SupportWeekInput) (*SupportWeek, error) {
	firstDay, err := DateToDay(input.FirstDay)
	if err != nil {
		return nil, err
	}

	lastDay, err := DateToDay(input.LastDay)
	if err != nil {
		return nil, err
	}

	return &SupportWeek{
		FirstDay: firstDay,
		LastDay:  lastDay,
		DevId:    input.DevId,
	}, nil
}

func newTask(input *TaskInput) (*Task, error) {
	attrs := make(map[DeveloperId]*Attribution, len(input.Attributions))

	for devId, input := range input.Attributions {
		attr, err := newAttribution(input)
		if err != nil {
			return nil, fmt.Errorf("error in creating task for %+v: %s", input, err)
		}
		attrs[devId] = attr
	}

	return &Task{
		Name:         input.Name,
		Attributions: attrs,
	}, nil
}

func newAttribution(input *AttributionInput) (*Attribution, error) {
	var firstDay *Day
	var lastDay *Day

	if input.FirstDay != nil {
		day, err := DateToDay(*input.FirstDay)
		if err != nil {
			return nil, err
		}
		firstDay = &day
	}

	if input.LastDay != nil {
		day, err := DateToDay(*input.LastDay)
		if err != nil {
			return nil, err
		}
		lastDay = &day
	}

	return &Attribution{
		EffortDays: input.Effort,
		FirstDay:   firstDay,
		LastDay:    lastDay,
	}, nil
}

func NewPlanningInput(planning *Planning) *PlanningInput {
	holidays := make([]string, len(planning.Holidays))
	for i, day := range planning.Holidays {
		holidays[i] = DayToDate(day)
	}

	developers := make([]*DeveloperInput, len(planning.Developers))
	for i, developer := range planning.Developers {
		offDays := make([]string, developer.OffDays.Len())
		for j, day := range developer.OffDays {
			offDays[j] = DayToDate(day)
		}
		var starts *string
		if developer.Starts != nil {
			date := DayToDate(*developer.Starts)
			starts = &date
		}
		var leaves *string
		if developer.Leaves != nil {
			date := DayToDate(*developer.Leaves)
			leaves = &date
		}

		developers[i] = &DeveloperInput{
			Id:      developer.Id,
			OffDays: offDays,
			Starts:  starts,
			Leaves:    leaves,
		}
	}

	supportWeeks := make([]*SupportWeekInput, len(planning.SupportWeeks))
	for i, week := range planning.SupportWeeks {
		supportWeeks[i] = &SupportWeekInput{
			FirstDay: DayToDate(week.FirstDay),
			LastDay:  DayToDate(week.LastDay),
			DevId:    week.DevId,
		}
	}

	tasks := make([]*TaskInput, len(planning.Tasks))
	for i, task := range planning.Tasks {
		attributions := make(map[DeveloperId]*AttributionInput)
		for developerId, attribution := range task.Attributions {
			attributions[developerId] = newAttributionInput(attribution)
		}

		tasks[i] = &TaskInput{
			Name:         task.Name,
			Attributions: attributions,
		}
	}

	return &PlanningInput{
		StartDay:     DayToDate(planning.StartDay),
		Holidays:     holidays,
		Developers:   developers,
		SupportWeeks: supportWeeks,
		Tasks:        tasks,
	}
}

func newAttributionInput(attr *Attribution) *AttributionInput {
	var firstDay *string
	if attr.FirstDay != nil {
		date := DayToDate(*attr.FirstDay)
		firstDay = &date
	}

	var lastDay *string
	if attr.LastDay != nil {
		date := DayToDate(*attr.LastDay)
		lastDay = &date
	}

	return &AttributionInput{
		Effort:   attr.EffortDays,
		FirstDay: firstDay,
		LastDay:  lastDay,
	}
}
