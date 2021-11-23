package state

import (
	"reflect"
	"testing"
)

var (
	testOn      State = "on"
	testOff     State = "off"
	testInvalid State = "foo"
	testOnEvent Event = func(s State) State {
		if s == testOn {
			return testOff
		} else {
			return testOn
		}
	}
	testOffEvent Event = func(s State) State {
		if s == testOff {
			return testOn
		} else {
			return testOff
		}
	}
)

func TestStateMachine_ProcessState(t *testing.T) {
	type args struct {
		state State
	}
	tests := []struct {
		name    string
		s       *StateMachine
		args    args
		want    State
		wantErr bool
	}{
		{
			name: "success on to off",
			s: &StateMachine{
				states: map[State]Event{
					testOn:  testOnEvent,
					testOff: testOffEvent,
				},
			},
			args: args{
				state: testOn,
			},
			want: testOff,
		},
		{
			name: "success off to on",
			s: &StateMachine{
				states: map[State]Event{
					testOn:  testOnEvent,
					testOff: testOffEvent,
				},
			},
			args: args{
				state: testOff,
			},
			want: testOn,
		},
		{
			name: "error invalid state",
			s: &StateMachine{
				states: map[State]Event{
					testOn:  testOnEvent,
					testOff: testOffEvent,
				},
			},
			args: args{
				state: testInvalid,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.ProcessState(tt.args.state)
			if (err != nil) != tt.wantErr {
				t.Errorf("StateMachine.ProcessState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StateMachine.ProcessState() = %v, want %v", got, tt.want)
			}
		})
	}
}
