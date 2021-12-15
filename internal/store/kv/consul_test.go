package kv

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/store/inf"
)

func TestConsul_crud(t *testing.T) {
	storeClient, err := NewConsul(context.Background(), "unittests")
	if err != nil {
		t.Errorf("cannot connect consul %v", err)
		return
	}

	type fields struct {
		client    inf.CRUD
		operation string
	}
	type args struct {
		key   string
		value []byte
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErr     bool
		wantErrCode int
		wantVal     [][]byte
	}{
		{
			name: "put hello",
			fields: fields{
				client:    storeClient,
				operation: "put",
			},
			args: args{
				key:   "hello",
				value: []byte("TEST mEssAgE"),
			},
		},
		{
			name: "get hello",
			fields: fields{
				client:    storeClient,
				operation: "get",
			},
			args: args{
				key: "hello",
			},
			wantVal: [][]byte{[]byte("TEST mEssAgE")},
		},
		{
			name: "delete hello",
			fields: fields{
				client:    storeClient,
				operation: "delete",
			},
			args: args{
				key: "hello",
			},
		},
		{
			name: "get hello again",
			fields: fields{
				client:    storeClient,
				operation: "get",
			},
			args: args{
				key: "hello",
			},
			wantErr:     true,
			wantErrCode: 404,
		},
		{
			name: "put some under hello",
			fields: fields{
				client:    storeClient,
				operation: "put",
			},
			args: args{
				key:   "hello/some",
				value: []byte("SOME mEssAgE"),
			},
		},
		{
			name: "get some under hello",
			fields: fields{
				client:    storeClient,
				operation: "get",
			},
			args: args{
				key: "hello/some",
			},
			wantVal: [][]byte{[]byte("SOME mEssAgE")},
		},
		{
			name: "delete hello/some",
			fields: fields{
				client:    storeClient,
				operation: "delete",
			},
			args: args{
				key: "hello/some",
			},
		},
		{
			name: "get unknown data",
			fields: fields{
				client:    storeClient,
				operation: "get",
			},
			args: args{
				key: "sdfsafdyokartikfdsf",
			},
			wantErr:     true,
			wantErrCode: 404,
		},
		{
			name: "post data x",
			fields: fields{
				client:    storeClient,
				operation: "post",
			},
			args: args{
				key:   "hello/x",
				value: []byte("SOME mEssAgE"),
			},
		},
		{
			name: "post data x again",
			fields: fields{
				client:    storeClient,
				operation: "post",
			},
			args: args{
				key:   "hello/x",
				value: []byte("SOME mEssAgE"),
			},
			wantErr:     true,
			wantErrCode: 409,
		},
		{
			name: "post data x again",
			fields: fields{
				client:    storeClient,
				operation: "post",
			},
			args: args{
				key:   "hello/x",
				value: []byte("SOME mEssAgE"),
			},
			wantErr:     true,
			wantErrCode: 409,
		},
		{
			name: "delete hello/x",
			fields: fields{
				client:    storeClient,
				operation: "delete",
			},
			args: args{
				key: "hello/x",
			},
			wantErr: false,
		},
		{
			name: "delete hello/x again",
			fields: fields{
				client:    storeClient,
				operation: "delete",
			},
			args: args{
				key: "hello/x",
			},
			wantErr:     true,
			wantErrCode: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields.client
			switch tt.fields.operation {
			case "put":
				if err := c.Put(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
					t.Errorf("Consul.%s() error = %v, wantErr %v", strings.Title(tt.fields.operation), err, tt.wantErr)
				} else if err != nil && tt.wantErrCode != err.GetCode() {
					t.Errorf("wantErrCode %v, but value is %v", tt.wantErrCode, err.GetCode())
				}
			case "post":
				if err := c.Post(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
					t.Errorf("Consul.%s() error = %v, wantErr %v", strings.Title(tt.fields.operation), err, tt.wantErr)
				} else if err != nil && tt.wantErrCode != err.GetCode() {
					t.Errorf("wantErrCode %v, but value is %v", tt.wantErrCode, err.GetCode())
				}
			case "get":
				if val, err := c.Get(tt.args.key); (err != nil) != tt.wantErr {
					t.Errorf("Consul.%s() error = %v, wantErr %v", strings.Title(tt.fields.operation), err, tt.wantErr)
				} else if err != nil && tt.wantErrCode != err.GetCode() {
					t.Errorf("wantErrCode %v, but value is %v", tt.wantErrCode, err.GetCode())
				} else {
					if len(val) != len(tt.wantVal) {
						t.Errorf("Consul.%s() val = %s, wantVal %s", strings.Title(tt.fields.operation), val, tt.wantVal)
					} else {
						for i := range val {
							if bytes.Compare(val[i], tt.wantVal[i]) != 0 {
								t.Errorf("Consul.%s() val = %s, wantVal %s", strings.Title(tt.fields.operation), val[i], tt.wantVal[i])
							}
						}
					}
				}
			case "delete":
				if err := c.Delete(tt.args.key); (err != nil) != tt.wantErr {
					t.Errorf("Consul.%s() error = %v, wantErr %v", strings.Title(tt.fields.operation), err, tt.wantErr)
				}
			default:
				t.Errorf("operation %v not found", tt.fields.operation)
			}
		})
	}
}
