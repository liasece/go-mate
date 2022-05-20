package writer

import (
	"testing"
)

func Test_getProtoFromStr(t *testing.T) {
	type args struct {
		originContent string
		typ           string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "OK",
			args: args{
				originContent: `
				message ItemStorage {
					repeated ItemStorageItem items = 1;
				}
				
				message User {
					string id = 1;
					string name = 2;
				}
				
				message UserUpdater {
					optional string name = 1;
				}
				`,
				typ: "UserUpdater",
			},
			want: `				message UserUpdater {
					optional string name = 1;
				}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getProtoFromStr(tt.args.originContent, tt.args.typ); got != tt.want {
				t.Errorf("getProtoFromStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_snakeString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "OK",
			args: args{
				s: "ID",
			},
			want: "id",
		},
		{
			name: "OK",
			args: args{
				s: "UserIDEq",
			},
			want: "user_id_eq",
		},
		{
			name: "OK",
			args: args{
				s: "aUserIDEqX",
			},
			want: "a_user_id_eq_x",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := snakeString(tt.args.s); got != tt.want {
				t.Errorf("snakeString() = %v, want %v", got, tt.want)
			}
		})
	}
}
