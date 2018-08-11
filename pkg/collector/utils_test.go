package collector

import (
	"testing"
)

func TestNewV3Client(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name       string
		args       args
		wantClient bool
		wantCtx    bool
	}{
		{
			name: "Anon",
			args: args{

			},
			wantClient: true,
			wantCtx:    true,
		}, {
			name: "Tokent",
			args: args{
				token: "12321321",
			},
			wantClient: true,
			wantCtx:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClient, gotCtx := NewV3Client(tt.args.token)
			if (gotClient != nil) != tt.wantClient {
				t.Errorf("NewV3Client() gotClient = %v, want %v", (gotClient != nil), tt.wantClient)

			}
			if (gotCtx != nil) != tt.wantCtx {
				t.Errorf("NewV3Client() gotCtx = %v, want %v", (gotCtx != nil), tt.wantCtx)

			}

		})
	}
}
