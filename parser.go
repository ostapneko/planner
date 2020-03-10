package planner

import "fmt"

type PlanningInput struct {
	Calendar []string
	Developers []*DeveloperInput `yaml:"developers"`
	SupportWeeks []*SupportWeekInput `yaml:"supportWeeks"`
	Tasks []*TaskInput `yaml:"tasks"`
}

type TaskInput struct {
	Name string
	Attributions map[DeveloperId]*AttributionInput
}

type AttributionInput struct {
	Duration DurationDays
	FirstDay *string `yaml:"firstDay"`
	LastDay *string `yaml:"lastDay"`
}

type SupportWeekInput struct {
	FirstDay string `yaml:"firstDay"`
	LastDay  string `yaml:"lastDay"`
	DevId    DeveloperId `yaml:"devId"`
}

type DeveloperInput struct {
	Id      DeveloperId
	OffDays []string `yaml:"offDays"`
}

func NewPlanning(input PlanningInput) (*Planning, error) {
	cal := make([]Day, len(input.Calendar))
	for i, s := range input.Calendar {
		d, err := DateToDay(s)
		if err != nil {
			return nil, fmt.Errorf("error parsing cal %s", err)
		}
		cal[i] = d
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

	return &Planning{
		Calendar:     cal,
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

	return &Developer{
		Id:      input.Id,
		OffDays: offDays,
	}, nil
}

func newSupportWeek(input *SupportWeekInput)(*SupportWeek, error) {
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

func newTask(input *TaskInput)(*Task, error) {
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

func newAttribution(input *AttributionInput)(*Attribution, error) {
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
		DurationDays: input.Duration,
		FirstDay:     firstDay,
		LastDay:      lastDay,
	}, nil
}
