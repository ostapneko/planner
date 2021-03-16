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
	var day5 Day = 5

	holidays := []Day{3, 4}

	attributions1 := map[DeveloperId]*Attribution{
		dev1Id: {
			EffortDays: 2,
			FirstDay:   &day1,
			LastDay:    &day5,
		},
	}

	task1 := &Task{
		Name:         "Task",
		Attributions: attributions1,
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
				Holidays:     holidays,
			}},
			wantErr: false,
		},
		{
			name: "Developer missing from Attributions",
			args: args{&Planning{
				Tasks:        []*Task{task1},
				Developers:   []*Developer{dev2},
				SupportWeeks: []*SupportWeek{sw2},
				Holidays:     holidays,
			}},
			wantErr: true,
		},
		{
			name: "support week invalid",
			args: args{&Planning{
				Tasks:        []*Task{task1},
				Developers:   []*Developer{dev1},
				SupportWeeks: []*SupportWeek{invalidSw},
				Holidays:     holidays,
			}},
			wantErr: true,
		},
		{
			name: "overlapping support week",
			args: args{&Planning{
				Tasks:        []*Task{task1},
				Developers:   []*Developer{dev1},
				SupportWeeks: []*SupportWeek{sw1, sw1},
				Holidays:     holidays,
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
			"dev3": {EffortDays: 1},
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
			"dev3": {EffortDays: 1},
		},
	}
	task3 := &Task{
		Name: "task3",
		Attributions: map[DeveloperId]*Attribution{
			"dev1": {EffortDays: 1},
			"dev2": {EffortDays: 10},
			"dev3": {EffortDays: 1},
		},
	}
	planning := &Planning{
		StartDay: 1,
		Holidays: []Day{9, 10},
		Developers: []*Developer{
			{
				Id:      "dev1",
				OffDays: []Day{1},
				Utilization: 1,
			},
			{
				Id:      "dev2",
				OffDays: []Day{2},
				Utilization: 1,
			},
			{
				Id:          "dev3",
				Utilization: 0.25,
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

	ForecastCompletion(planning)

	// start day: 1
	// holidays: 9, 10
	// support week: dev1 @ 7, 8
	// dev1 off days: 1
	// dev2 off days: 2
	// saturdays  + sundays: 2, 3, 9, 10, 16, 17
	// task1:
	// dev1 (3d) : 4, 5, 6
	// dev2 (2d): 1, 4
	// dev3 (1d / 0.25 = 4d): 1, 4, 5, 6
	// last day: 6
	// task2:
	// dev1 (3d): 11, 12, 13
	// dev2 (2d): 5, 6
	// dev3 (1d / 0.25 = 4d): 7, 8, 11, 12
	// last day: 13
	// task3:
	// dev1 (1d): 14
	// dev2 (10d): 7, 8, 11, 12, 13, 14, 15, 18, 19, 20
	// dev3 (1d / 0.25 = 4d): 13, 14, 15, 18
	// last day: 20

	examples := []struct {
		act Day
		exp Day
	}{
		{act: *task1.Attributions["dev1"].FirstDay, exp: 4},
		{act: *task1.Attributions["dev1"].LastDay, exp: 6},
		{act: *task1.Attributions["dev2"].FirstDay, exp: 1},
		{act: *task1.Attributions["dev2"].LastDay, exp: 4},
		{act: *task1.Attributions["dev3"].FirstDay, exp: 1},
		{act: *task1.Attributions["dev3"].LastDay, exp: 6},
		{act: *task1.LastDay, exp: 6},
		{act: *task2.Attributions["dev1"].FirstDay, exp: 11},
		{act: *task2.Attributions["dev1"].LastDay, exp: 13},
		{act: *task2.Attributions["dev2"].FirstDay, exp: 5},
		{act: *task2.Attributions["dev2"].LastDay, exp: 6},
		{act: *task2.Attributions["dev3"].FirstDay, exp: 7},
		{act: *task2.Attributions["dev3"].LastDay, exp: 12},
		{act: *task2.LastDay, exp: 13},
		{act: *task3.Attributions["dev1"].FirstDay, exp: 14},
		{act: *task3.Attributions["dev1"].LastDay, exp: 14},
		{act: *task3.Attributions["dev2"].FirstDay, exp: 7},
		{act: *task3.Attributions["dev2"].LastDay, exp: 20},
		{act: *task3.Attributions["dev3"].FirstDay, exp: 13},
		{act: *task3.Attributions["dev3"].LastDay, exp: 18},
		{act: *task3.LastDay, exp: 20},
	}

	for i, example := range examples {
		if example.act != example.exp {
			t.Errorf("exp %d, got %d in example %d", example.exp, example.act, i+1)
		} else {
			t.Logf("pass example %d", i)
		}
	}
}
