package transfer

import (
	"testing"

	"github.com/go-test/deep"
)

func TestDataToBytes(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name string
		args args
		want []byte
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
				data: map[string]interface{}{
					"key": "value",
				},
			},
			want: []byte(`{"key":"value"}`),
		},
		{
			name: "[]interface{}",
			args: args{
				data: []interface{}{
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
			want: []byte(`["string",1,1.1,true,null,{"key":"value"}]`),
		},
		{
			name: "[]byte",
			args: args{
				data: []byte("byte"),
			},
			want: []byte("byte"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DataToBytes(tt.args.data)
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("DataToBytes() = %v", diff)
			}
		})
	}
}
