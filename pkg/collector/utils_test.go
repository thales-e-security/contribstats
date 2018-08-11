package collector

import (
	"github.com/spf13/viper"
	"github.com/thales-e-security/contribstats/pkg/config"
	"testing"
)

func init() {
	config.InitConfig("")
}

func TestNewV3Client(t *testing.T) {

	tests := []struct {
		name       string
		wantClient bool
		wantAuth   bool
		wantCtx    bool
	}{
		{
			name:       "Anon",
			wantClient: true,
			wantAuth:   false,
			wantCtx:    true,
		}, {
			name:       "Token",
			wantClient: true,
			wantAuth:   true,
			wantCtx:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantAuth {
				viper.Set("token", nil)
			}
			gotClient, gotCtx := NewV3Client()
			if (gotClient != nil) != tt.wantClient {
				t.Errorf("NewV3Client() gotClient = %v, want %v", (gotClient != nil), tt.wantClient)

			}
			if (gotCtx != nil) != tt.wantCtx {
				t.Errorf("NewV3Client() gotCtx = %v, want %v", (gotCtx != nil), tt.wantCtx)

			}

		})
	}
}
