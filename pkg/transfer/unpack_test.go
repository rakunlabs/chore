package transfer

import (
	"testing"

	"github.com/go-test/deep"
)

func TestBytesToData(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "nil",
			args: args{
				data: nil,
			},
			want: nil,
		},
		{
			name: "map[string]interface{}",
			args: args{
				data: []byte(`{"key":"value"}`),
			},
			want: map[string]interface{}{
				"key": "value",
			},
		},
		{
			name: "[]interface{}",
			args: args{
				data: []byte(`["string",1,1.1,true,null,{"key":"value"}]`),
			},
			want: []interface{}{
				"string",
				1,
				1.1,
				true,
				nil,
				map[string]interface{}{
					"key": "value",
				},
			},
		},
		{
			name: "[]byte",
			args: args{
				data: []byte("byte"),
			},
			want: string("byte"),
		},
		{
			name: "string",
			args: args{
				data: []byte("string"),
			},
			want: "string",
		},
		{
			name: "int",
			args: args{
				data: []byte("1"),
			},
			want: 1,
		},
		{
			name: "float",
			args: args{
				data: []byte("1.1"),
			},
			want: 1.1,
		},
		{
			name: "bool",
			args: args{
				data: []byte("true"),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BytesToData(tt.args.data)
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("BytesToData() = %v", diff)
			}
		})
	}
}
