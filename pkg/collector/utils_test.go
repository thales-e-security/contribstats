package collector

import (
	"github.com/thales-e-security/contribstats/pkg/config"
	"log"
	"testing"
)

func init() {
	config.InitConfig("")
}

func ExampleNewV3Client() {

	c, _ := NewV3Client(config.Constants{})
	log.Println(c)
}

func TestNewV3Client(t *testing.T) {

	tests := []struct {
		name       string
		wantClient bool
		wantAuth   bool
		wantCtx    bool
		constants  config.Constants
	}{
		{
			name:       "Anon",
			wantClient: true,
			wantAuth:   false,
			wantCtx:    true,
			constants:  config.Constants{},
		}, {
			name:       "Token",
			wantClient: true,
			wantAuth:   true,
			wantCtx:    true,
			constants:  config.Constants{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantAuth {
				tt.constants.Token = ""
			}
			gotClient, gotCtx := NewV3Client(tt.constants)
			if (gotClient != nil) != tt.wantClient {
				t.Errorf("NewV3Client() gotClient = %v, want %v", (gotClient != nil), tt.wantClient)

			}
			if (gotCtx != nil) != tt.wantCtx {
				t.Errorf("NewV3Client() gotCtx = %v, want %v", (gotCtx != nil), tt.wantCtx)

			}

		})
	}
}
