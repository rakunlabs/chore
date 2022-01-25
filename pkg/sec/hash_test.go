package sec

import (
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {
	type args struct {
		password string
		hash     string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "generate and check",
			args: args{
				password: "mypassword",
			},
			want: true,
		},
		{
			name: "generate and check empty",
			args: args{
				password: "",
			},
			want: true,
		},
		{
			name: "generate and check wrongly",
			args: args{
				password: "mypassword",
				hash:     "dummyhash",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.hash == "" {
				var err error
				tt.args.hash, err = HashPassword(tt.args.password)
				if err != nil {
					t.Errorf("HashPassword; %v", err)
					return
				}
			}

			// fmt.Println(tt.args.hash)

			if got := CheckHashPassword(tt.args.hash, tt.args.password); got != tt.want {
				t.Errorf("CheckPasswordHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
