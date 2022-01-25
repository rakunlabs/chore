package sec

import (
	"reflect"
	"testing"
	"time"
)

func TestJWT_Generate(t *testing.T) {
	type fields struct {
		secret []byte
	}
	type args struct {
		claims  map[string]interface{}
		expDate int64
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantValidateErr bool
		wantErr         bool
	}{
		{
			name: "simple false",
			fields: fields{
				secret: []byte("pass1234"),
			},
			args:            args{},
			wantValidateErr: false,
			wantErr:         false,
		},
		{
			name: "with 1 hour",
			fields: fields{
				secret: []byte("pass1234"),
			},
			args: args{
				claims:  map[string]interface{}{"info": "hello"},
				expDate: time.Now().Add(time.Hour).Unix(),
			},
			wantValidateErr: false,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &JWT{
				secret: tt.fields.secret,
			}
			got, err := tr.Generate(tt.args.claims, tt.args.expDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("JWT.Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			claims, err := tr.Validate(got)
			if (err != nil) != tt.wantValidateErr {
				t.Errorf("JWT.Validate() error = %v, wantValidateErr %v", err, tt.wantValidateErr)
				return
			}

			delete(claims, "exp")
			if len(tt.args.claims) != 0 && !reflect.DeepEqual(tt.args.claims, claims) {
				t.Errorf("claims %+v, want %+v", claims, tt.args.claims)
			}
		})
	}
}
