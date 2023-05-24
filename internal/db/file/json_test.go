package file

import (
	"testing"

	"github.com/brandonlbarrow/gonk/v2/internal/db"
	"github.com/go-test/deep"
)

func TestNewJSONFileDB(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *JSONFileDB
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				path: "../../../db.json",
			},
			want: &JSONFileDB{
				db: db.ServerMap{
					"1": db.Server{
						GuildID: "1",
						Streamers: []db.Streamer{
							{
								UserID:    "2",
								ChannelID: "3",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewJSONFileDB(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewJSONFileDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(deep.Equal(got.db, tt.want.db)) > 0 {
				t.Error(deep.Equal(got.db, tt.want.db))
			}
		})
	}
}
