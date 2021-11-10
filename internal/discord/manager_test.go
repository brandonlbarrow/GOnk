package discord

import (
	"reflect"
	"testing"
)

func TestNewManager(t *testing.T) {
	type args struct {
		opts []ManagerOption
	}
	tests := []struct {
		name string
		args args
		want *Manager
	}{
		{
			name: "success",
			args: args{
				opts: []ManagerOption{},
			},
			want: &Manager{},
		},
		{
			name: "success with option",
			args: args{
				opts: []ManagerOption{
					WithGuildID("foo"),
				},
			},
			want: &Manager{guildID: "foo"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewManager(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewManager() = %v, want %v", got, tt.want)
			}
		})
	}
}
