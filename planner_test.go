package planner

import "testing"

func Test_checkGraph(t *testing.T) {
	type args struct {
		planning *Planning
	}

	dev1Id := DeveloperId("dev1")

	var day1 Day = 1
	var day2 Day = 2
	var day5 Day = 5

	attributions1 := map[DeveloperId]*Attribution{
		dev1Id: {
			DurationDays: 5,
			FirstDay:     &day1,
			LastDay:      &day5,
		},
	}

	attributionsWrongDuration := map[DeveloperId]*Attribution{
		dev1Id: {
			DurationDays: 5,
			FirstDay:     &day2,
			LastDay:      &day5,
		},
	}

	attributionsOnlyLastDay := map[DeveloperId]*Attribution{
		dev1Id: {
			DurationDays: 5,
			LastDay:      &day5,
		},
	}

	task1 := &Task{
		Name:         "Task",
		Attributions: attributions1,
	}

	taskAttributionWrongDuration := &Task{
		Name:         "WrongTask",
		Attributions: attributionsWrongDuration,
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

	dev1 := &Developer{
		Id:      dev1Id,
		OffDays: []Day{},
	}

	dev2 := &Developer{
		Id:      "dev2",
		OffDays: []Day{},
	}

	sw1 := &SupportWeek{
		FirstDay: 10,
		LastDay:  12,
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
				Calendar:     []Day{},
			}},
			wantErr: false,
		},
		{
			name: "Developer missing from Attributions",
			args: args{&Planning{
				Tasks:        []*Task{task1},
				Developers:   []*Developer{dev2},
				SupportWeeks: []*SupportWeek{sw2},
				Calendar:     []Day{},
			}},
			wantErr: true,
		},
		{
			name: "Wrong attribution duration",
			args: args{&Planning{
				Tasks:        []*Task{taskAttributionWrongDuration},
				Developers:   []*Developer{dev1},
				SupportWeeks: []*SupportWeek{sw1},
				Calendar:     []Day{},
			}},
			wantErr: true,
		},
		{
			name: "Wrong attribution duration",
			args: args{&Planning{
				Tasks:        []*Task{taskAttributionOnlyLastDay},
				Developers:   []*Developer{dev1},
				SupportWeeks: []*SupportWeek{sw1},
				Calendar:     []Day{},
			}},
			wantErr: true,
		},
		{
			name: "Task last day but no attribution last day",
			args: args{&Planning{
				Tasks:        []*Task{taskLastDayNoAttribution},
				Developers:   []*Developer{dev1},
				SupportWeeks: []*SupportWeek{sw1},
				Calendar:     []Day{},
			}},
			wantErr: true,
		},
		{
			name: "support week invalid",
			args: args{&Planning{
				Tasks:        []*Task{task1},
				Developers:   []*Developer{dev1},
				SupportWeeks: []*SupportWeek{invalidSw},
				Calendar:     []Day{},
			}},
			wantErr: true,
		},
		{
			name: "overlapping support week",
			args: args{&Planning{
				Tasks:        []*Task{task1},
				Developers:   []*Developer{dev1},
				SupportWeeks: []*SupportWeek{sw1, sw1},
				Calendar:     []Day{},
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
