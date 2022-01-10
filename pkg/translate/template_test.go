package translate

import (
	"reflect"
	"testing"
)

func TestTemplate_Ext(t *testing.T) {
	tr := NewTemplate()

	type args struct {
		v    map[string]interface{}
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "basic test",
			args: args{
				v: map[string]interface{}{
					"name": "chore",
				},
				file: "my name is {{.name}}",
			},
			want:    []byte("my name is chore"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tr.Ext(tt.args.v, tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("Template.Ext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Template.Ext() = %v, want %v", got, tt.want)
			}
		})
	}
}
