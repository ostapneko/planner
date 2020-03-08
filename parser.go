package planner

import "fmt"

type PlanningInput struct {
	Calendar []string
	Developers []*DeveloperInput `yaml:"developers"`
	SupportWeeks []*SupportWeekInput `yaml:"supportWeeks"`
	Tasks []*Task `yaml:"tasks"`
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
		dev, err := NewDeveloper(input)
		if err != nil {
			return nil, fmt.Errorf("error parsing developer %s", err)
		}
		devs[i] = dev
	}

	weeks := make([]*SupportWeek, len(input.SupportWeeks))
	for i, input := range input.SupportWeeks {
		week, err := NewSupportWeek(input)
		if err != nil {
			return nil, fmt.Errorf("error parsing week")
		}
		weeks[i] = week
	}

	return &Planning{
		Calendar:     cal,
		Developers:   devs,
		SupportWeeks: weeks,
		Tasks:        input.Tasks,
	}, nil
}

func NewDeveloper(input *DeveloperInput) (*Developer, error) {
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

func NewSupportWeek(input *SupportWeekInput)(*SupportWeek, error) {
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
