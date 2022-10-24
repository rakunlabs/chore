package js

import (
	"context"
	"errors"
	"testing"

	"github.com/go-test/deep"
)

func TestGoja_RunScript(t *testing.T) {
	type args struct {
		script string
		inputs []interface{}
	}
	tests := []struct {
		name        string
		args        args
		function    func(x *interface{}) func(interface{})
		functionRes interface{}
		want        []byte
		wantErr     bool
		wantErrType error
	}{
		{
			name: "test basic",
			args: args{
				script: `
				function main(v1, v2) {
					return v1+v2
				}
				`,
				inputs: []interface{}{
					1, 4,
				},
			},
			want:    []byte("5"),
			wantErr: false,
		},
		{
			name: "test basic",
			args: args{
				script: `
				function main(v1) {
					return v1
				}
				`,
				inputs: []interface{}{
					[]byte("hello"),
				},
			},
			want:    []byte("hello"),
			wantErr: false,
		},
		{
			name: "reference error",
			args: args{
				script: `
				function main(v1) {
					return v2;
				}
				`,
			},
			want:    []byte("ReferenceError: v2 is not defined at main (<eval>:3:13(1))"),
			wantErr: true,
		},
		{
			name: "throw error",
			args: args{
				script: `
				function main(v1) {
					throw "upps";
					return v1;
				}
				`,
			},
			want:        []byte("upps"),
			wantErr:     true,
			wantErrType: ErrThrow,
		},
		{
			name: "syntax error",
			args: args{
				script: `
				function main(v1) {
					return v1,,
				}
				`,
			},
			wantErr: true,
		},
		{
			name: "additional functions",
			args: args{
				script: `
				function main(v1) {
					checkValue(v1);
					return v1;
				}
				`,
				inputs: []interface{}{"hello"},
			},
			functionRes: "hello",
			function: func(x *interface{}) func(interface{}) {
				return func(v interface{}) {
					*x = v
				}
			},
			want:    []byte("hello"),
			wantErr: false,
		},
		{
			name: "additional functions object",
			args: args{
				script: `
				function main() {
					checkValue({"hello": "world"});
				}
				`,
			},
			functionRes: map[string]interface{}{
				"hello": "world",
			},
			function: func(x *interface{}) func(interface{}) {
				return func(v interface{}) {
					*x = v
				}
			},
			wantErr: false,
		},
	}

	g := NewGoja()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var checkValue interface{}
			if tt.function != nil {
				fn := tt.function(&checkValue)
				g.SetFunction("checkValue", fn)
			}

			got, err := g.RunScript(context.Background(), tt.args.script, tt.args.inputs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Goja.RunScript() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && tt.wantErrType != nil && !errors.Is(err, tt.wantErrType) {
				t.Errorf("Goja.RunScript() error = %v, wantErrType %v", err, tt.wantErrType)
			}
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("Goja.RunScript() = %s, got=%s, want=%s", diff, got, tt.want)
			}

			if tt.function != nil {
				if diff := deep.Equal(checkValue, tt.functionRes); diff != nil {
					t.Errorf("Goja.RunScript() = %s, got=%s, want=%s", diff, checkValue, tt.functionRes)
				}
			}
		})
	}
}
