package planner

import (
	"testing"
)

func Test_checkPlanning(t *testing.T) {
	type args struct {
		planning *Planning
	}

	dev1Id := DeveloperId("dev1")

	var day1 Day = 1
	var day2 Day = 2
	var day5 Day = 5

	calendar := []Day{1, 2, 3, 4}

	attributions1 := map[DeveloperId]*Attribution{
		dev1Id: {
			EffortDays: 2,
			FirstDay:   &day1,
			LastDay:    &day5,
		},
	}

	attributionsWrongEffort := map[DeveloperId]*Attribution{
		dev1Id: {
			EffortDays: 5,
			FirstDay:   &day2,
			LastDay:    &day5,
		},
	}

	attributionsOnlyLastDay := map[DeveloperId]*Attribution{
		dev1Id: {
			EffortDays: 5,
			LastDay:    &day5,
		},
	}

	task1 := &Task{
		Name:         "Task",
		Attributions: attributions1,
	}

	taskAttributionWrongEffort := &Task{
		Name:         "WrongTask",
		Attributions: attributionsWrongEffort,
	}

	taskAttributionOnlyLastDay := &Task{
		Name:         "Only Last Day",
		Attributions: attributionsOnlyLastDay,
	}

	taskLastDayNoAttribution := &Task{
		Name:         "No attribution with last day",
		Attributions: make(map[DeveloperId]*Attribution, 0),
		LastDay:      &day1,
	}

	taskInconsistentLastDay := &Task{
		Name:         "Inconsistent last day",
		Attributions: make(map[DeveloperId]*Attribution, day2),
		LastDay:      &day1,
	}

	dev1 := &Developer{
		Id:      dev1Id,
		OffDays: []Day{3},
	}

	dev2 := &Developer{
		Id:      "dev2",
		OffDays: []Day{},
	}

	sw1 := &SupportWeek{
		FirstDay: 2,
		LastDay:  2,
		DevId:    dev1Id,
	}

	sw2 := &SupportWeek{
		FirstDay: 13,
		LastDay:  14,
		DevId:    "dev2",
	}

	invalidSw := &SupportWeek{
		FirstDay: 14,
		LastDay:  13,
		DevId:    dev1Id,
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid graph",
			args: args{&Planning{
				Tasks:        []*Task{task1},
				Developers:   []*Developer{dev1},
				SupportWeeks: []*SupportWeek{sw1},
				Calendar:     calendar,
			}},
			wantErr: false,
		},
		{
			name: "Developer missing from Attributions",
			args: args{&Planning{
				Tasks:        []*Task{task1},
				Developers:   []*Developer{dev2},
				SupportWeeks: []*SupportWeek{sw2},
				Calendar:     calendar,
			}},
			wantErr: true,
		},
		{
			name: "Wrong attribution effort",
			args: args{&Planning{
				Tasks:        []*Task{taskAttributionWrongEffort},
				Developers:   []*Developer{dev1},
				SupportWeeks: []*SupportWeek{sw1},
				Calendar:     calendar,
			}},
			wantErr: true,
		},
		{
			name: "Wrong attribution effort",
			args: args{&Planning{
				Tasks:        []*Task{taskAttributionOnlyLastDay},
				Developers:   []*Developer{dev1},
				SupportWeeks: []*SupportWeek{sw1},
				Calendar:     calendar,
			}},
			wantErr: true,
		},
		{
			name: "Task last day but no attribution last day",
			args: args{&Planning{
				Tasks:        []*Task{taskLastDayNoAttribution},
				Developers:   []*Developer{dev1},
				SupportWeeks: []*SupportWeek{sw1},
				Calendar:     calendar,
			}},
			wantErr: true,
		},
		{
			name: "Inconsistent last day",
			args: args{&Planning{
				Tasks:        []*Task{taskInconsistentLastDay},
				Developers:   []*Developer{dev1},
				SupportWeeks: []*SupportWeek{sw1},
				Calendar:     calendar,
			}},
			wantErr: true,
		},
		{
			name: "support week invalid",
			args: args{&Planning{
				Tasks:        []*Task{task1},
				Developers:   []*Developer{dev1},
				SupportWeeks: []*SupportWeek{invalidSw},
				Calendar:     calendar,
			}},
			wantErr: true,
		},
		{
			name: "overlapping support week",
			args: args{&Planning{
				Tasks:        []*Task{task1},
				Developers:   []*Developer{dev1},
				SupportWeeks: []*SupportWeek{sw1, sw1},
				Calendar:     calendar,
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckPlanning(tt.args.planning); (err != nil) != tt.wantErr {
				t.Errorf("CheckPlanning() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestForecastCompletion(t *testing.T) {
	var daySix Day = 6
	task1 := &Task{
		Name: "task1",
		Attributions: map[DeveloperId]*Attribution{
			"dev1": {EffortDays: 3},
			"dev2": {EffortDays: 2},
		},
	}
	task2 := &Task{
		Name: "task2",
		Attributions: map[DeveloperId]*Attribution{
			"dev1": {
				EffortDays: 3,
				FirstDay:   &daySix,
			},
			"dev2": {EffortDays: 2},
		},
	}
	task3 := &Task{
		Name: "task3",
		Attributions: map[DeveloperId]*Attribution{
			"dev1": {EffortDays: 1},
			"dev2": {EffortDays: 10},
		},
	}
	planning := &Planning{
		Calendar: []Day{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Developers: []*Developer{
			{
				Id:      "dev1",
				OffDays: []Day{1},
			},
			{
				Id:      "dev2",
				OffDays: []Day{2},
			},
		},
		SupportWeeks: []*SupportWeek{
			{
				FirstDay: 7,
				LastDay:  8,
				DevId:    "dev1",
			},
		},
		Tasks: []*Task{
			task1,
			task2,
			task3,
		},
	}

	ForecastCompletion(planning, 1)

	// task1:
	// dev1: 2, 3, 4
	// dev2: 1, 3
	// first day: 1
	// last day: 4
	// task2:
	// dev1: 5, 6, 9
	// dev2: 4, 5
	// first day: 4
	// last day: 9
	// task3:
	// dev1: 10
	// first day: 10
	// last day: 10
	// dev2: 6, 7, 8, 9, 10 (cannot be completed)
	// first day: 6
	// last day: 10

	examples := map[Day]Day{
		*task1.Attributions["dev1"].FirstDay: 2,
		*task1.Attributions["dev1"].LastDay:  4,
		*task1.Attributions["dev2"].FirstDay: 1,
		*task1.Attributions["dev2"].LastDay:  3,
		*task1.LastDay:                       4,
		*task2.Attributions["dev1"].FirstDay: 5,
		*task2.Attributions["dev1"].LastDay:  9,
		*task2.Attributions["dev2"].FirstDay: 4,
		*task2.Attributions["dev2"].LastDay:  5,
		*task2.LastDay:                       9,
		*task3.Attributions["dev1"].FirstDay: 10,
		*task3.Attributions["dev1"].LastDay:  10,
		*task3.Attributions["dev2"].FirstDay: 6,
	}

	i := 0
	for act, exp := range examples {
		i++
		if act != exp {
			t.Errorf("exp %d, got %d", exp, act)
		} else {
			t.Logf("pass example %d", i)
		}
	}

	if task3.Attributions["dev2"].LastDay != nil {
		t.Errorf("expected task3 attr2's last day to be nil, but found %d", *task3.Attributions["dev2"].LastDay)
	}

	if task3.LastDay != nil {
		t.Errorf("expected task3's last day to be nil but found %d", *task3.LastDay)
	}
}
