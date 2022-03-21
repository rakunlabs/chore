package email

import (
	"testing"
)

func TestClient_Send(t *testing.T) {
	type args struct {
		msg     []byte
		headers map[string][]string
	}
	tests := []struct {
		name    string
		client  Client
		args    args
		wantErr bool
	}{
		{
			name:   "test",
			client: NewClient("smtp.office365.com", 587, "eray.ates@ingenico.com", "---"),
			args: args{
				msg: []byte("<h1>this is test</h1>"),
				headers: map[string][]string{
					"Subject": {"test 2"},
					"From":    {"eray.ates@ingenico.com"},
					"To":      {"eray.ates@ingenico.com"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// SKIP THIS TEST
			t.SkipNow()
			if err := tt.client.Send(tt.args.msg, tt.args.headers); (err != nil) != tt.wantErr {
				t.Errorf("Client.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
