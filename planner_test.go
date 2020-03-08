package planner

import "testing"

func Test_checkGraph(t *testing.T) {
	type args struct {
		tasks        []task
		developers   []developer
		supportWeeks []supportWeek
		cal          calendar
	}
	task1 := task{
		name:         "task",
		attributions: map[developerId]durationDays{"dev1": 3},
	}

	dev1 := developer{
		id:    "dev1",
		offDays: []day{},
	}

	dev2 := developer{
		id:    "dev2",
		offDays: []day{},
	}

	sw1 := supportWeek{
		firstDay: 10,
		lastDay:  12,
		devId: "dev1",
	}

	sw2 := supportWeek{
		firstDay: 13,
		lastDay:  14,
		devId: "dev2",
	}

	invalidSw := supportWeek{
		firstDay: 14,
		lastDay:  13,
		devId: "dev1",
	}

	cal1 := calendar{
		days: []day{1, 2, 3},
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid graph",
			args: args{
				tasks:        []task{task1},
				developers:   []developer{dev1},
				supportWeeks: []supportWeek{sw1},
				cal:          cal1,
			},
			wantErr: false,
		},
		{
			name: "developer missing from attributions",
			args: args{
				tasks:        []task{task1},
				developers:   []developer{dev2},
				supportWeeks: []supportWeek{sw2},
				cal:          cal1,
			},
			wantErr: true,
		},
		{
			name: "support week invalid",
			args: args{
				tasks:        []task{task1},
				developers:   []developer{dev1},
				supportWeeks: []supportWeek{invalidSw},
				cal:          cal1,
			},
			wantErr: true,
		},
		{
			name: "overlapping support week",
			args: args{
				tasks:        []task{task1},
				developers:   []developer{dev1},
				supportWeeks: []supportWeek{sw1, sw1},
				cal:          cal1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkGraph(tt.args.tasks, tt.args.developers, tt.args.supportWeeks, tt.args.cal); (err != nil) != tt.wantErr {
				t.Errorf("checkGraph() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
