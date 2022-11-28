package utils

import (
	"reflect"
	"testing"
)

func TestTmplUtilsFunc_SplitCamelCase(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				str: "",
			},
			want: []string{},
		},
		{
			name: "test2",
			args: args{
				str: "Hello",
			},
			want: []string{"Hello"},
		},
		{
			name: "test3",
			args: args{
				str: "HelloWorld",
			},
			want: []string{"Hello", "World"},
		},
		{
			name: "test4",
			args: args{
				str: "idHelloWorld",
			},
			want: []string{"id", "Hello", "World"},
		},
		{
			name: "test5",
			args: args{
				str: "HelloIDWorld",
			},
			want: []string{"Hello", "ID", "World"},
		},
		{
			name: "test6",
			args: args{
				str: "iD",
			},
			want: []string{"i", "D"},
		},
		{
			name: "test7",
			args: args{
				str: "iPhone",
			},
			want: []string{"i", "Phone"},
		},
		{
			name: "test7",
			args: args{
				str: "[]*iPhone",
			},
			want: []string{"[]*", "i", "Phone"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SplitCamelCase(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TmplUtilsFunc.SplitCamelCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
