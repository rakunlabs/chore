package request

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Send(t *testing.T) {
	c := NewClient()

	type args struct {
		URL     string
		method  string
		payload []byte
	}
	tests := []struct {
		name    string
		server  func(*testing.T) *httptest.Server
		args    args
		wantErr bool
	}{
		{
			name: "request GET",
			server: func(t *testing.T) *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Method != http.MethodGet {
						t.Errorf("Method is %v, want %v", r.Method, http.MethodGet)
					}
				}))
			},
			args: args{
				method:  http.MethodGet,
				payload: nil,
			},
			wantErr: false,
		},
		{
			name: "request POST",
			server: func(t *testing.T) *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Method != http.MethodPost {
						t.Errorf("Method is %v, want %v", r.Method, http.MethodGet)
					}

					data, err := io.ReadAll(r.Body)
					if err != nil {
						t.Errorf("Error appered %v", err)
					}

					if bytes.Compare(data, []byte(`{"name":"test"}`)) != 0 {
						t.Errorf("data is %s, want %s", data, []byte(`{"name":"test"}`))
					}
				}))
			},
			args: args{
				method:  http.MethodPost,
				payload: []byte(`{"name":"test"}`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.server(t)
			if err := c.Send(context.TODO(), srv.URL, tt.args.method, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("Client.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
