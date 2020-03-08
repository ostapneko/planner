package planner

import "testing"

func Test_checkGraph(t *testing.T) {
	type args struct {
		tasks        []*Task
		developers   []*Developer
		supportWeeks []*SupportWeek
		cal          []Day
	}
	task1 := &Task{
		Name:         "Task",
		Attributions: map[DeveloperId]DurationDays{"dev1": 3},
	}

	dev1 := &Developer{
		Id:      "dev1",
		OffDays: []Day{},
	}

	dev2 := &Developer{
		Id:      "dev2",
		OffDays: []Day{},
	}

	sw1 := &SupportWeek{
		FirstDay: 10,
		LastDay:  12,
		DevId:    "dev1",
	}

	sw2 := &SupportWeek{
		FirstDay: 13,
		LastDay:  14,
		DevId:    "dev2",
	}

	invalidSw := &SupportWeek{
		FirstDay: 14,
		LastDay:  13,
		DevId:    "dev1",
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid graph",
			args: args{
				tasks:        []*Task{task1},
				developers:   []*Developer{dev1},
				supportWeeks: []*SupportWeek{sw1},
				cal:          []Day{},
			},
			wantErr: false,
		},
		{
			name: "Developer missing from Attributions",
			args: args{
				tasks:        []*Task{task1},
				developers:   []*Developer{dev2},
				supportWeeks: []*SupportWeek{sw2},
				cal:          []Day{},
			},
			wantErr: true,
		},
		{
			name: "support week invalid",
			args: args{
				tasks:        []*Task{task1},
				developers:   []*Developer{dev1},
				supportWeeks: []*SupportWeek{invalidSw},
				cal:          []Day{},
			},
			wantErr: true,
		},
		{
			name: "overlapping support week",
			args: args{
				tasks:        []*Task{task1},
				developers:   []*Developer{dev1},
				supportWeeks: []*SupportWeek{sw1, sw1},
				cal:          []Day{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckGraph(tt.args.tasks, tt.args.developers, tt.args.supportWeeks, tt.args.cal); (err != nil) != tt.wantErr {
				t.Errorf("CheckGraph() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
