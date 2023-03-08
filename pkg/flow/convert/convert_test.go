package convert

import (
	"reflect"
	"testing"
)

func TestGetList(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "string",
			args: args{
				value: "a b c",
			},
			want: []string{"a", "b", "c"},
		},
		{
			name: "string with comma",
			args: args{
				value: "a,b, c",
			},
			want: []string{"a", "b", "c"},
		},
		{
			name: "other type",
			args: args{
				value: 1,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetList(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBoolean(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "string",
			args: args{
				value: "true",
			},
			want: true,
		},
		{
			name: "string x",
			args: args{
				value: "x",
			},
			want: false,
		},
		{
			name: "bool",
			args: args{
				value: true,
			},
			want: true,
		},
		{
			name: "other type",
			args: args{
				value: 1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetBoolean(tt.args.value); got != tt.want {
				t.Errorf("GetBoolean() = %v, want %v", got, tt.want)
			}
		})
	}
}
