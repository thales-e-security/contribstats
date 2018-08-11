package server

import (
	"reflect"
	"testing"

	"github.com/spf13/viper"
	"net/http"
	"encoding/json"
	"net/http/httptest"
	"github.com/thales-e-security/contribstats/pkg/collector"
	"time"
)

func init() {
	viper.Set("interval", 60)
}

func TestNewStatServer(t *testing.T) {
	type args struct {
		debug bool
	}
	tests := []struct {
		name   string
		args   args
		wantSs *StatServer
	}{
		{
			name: "OK",
			args: args{
				debug: true,
			},
			wantSs: &StatServer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSs := NewStatServer(tt.args.debug); !reflect.DeepEqual(gotSs, tt.wantSs) {
				t.Errorf("NewStatServer() = %v, want %v", gotSs, tt.wantSs)
			}
		})
	}
}

func TestStatServer_Start(t *testing.T) {
	tests := []struct {
		name    string
		ss      *StatServer
		wantErr bool
	}{
		{
			name:    "OK",
			ss:      &StatServer{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// store old os.Exit
			oldOSExit := osExit
			// override os.Exit
			var got int
			osExit = func(code int) {
				got = code
			}
			// Start the server
			go tt.ss.Start()
			// wait for it...
			time.Sleep(100 * time.Millisecond)
			// Kill it ...
			cancel <- struct{}{}
			// repair os.Exit
			osExit = oldOSExit
			// See what we got
			if got != 0 {
				t.Error("Got unhealthy exit")
			}
		})
	}
}

func TestStatServer_startServer(t *testing.T) {
	type args struct {
		errs chan error
	}
	tests := []struct {
		name string
		ss   *StatServer
		args args
	}{

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ss.startServer(tt.args.errs)
		})
	}
}

func TestStatServer_startCollector(t *testing.T) {
	type args struct {
		errs chan error
	}
	tests := []struct {
		name string
		ss   *StatServer
		args args
	}{
		{
			name: "ok",
			ss:   &StatServer{},
			args: args{
				errs: errs,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ss.startCollector(tt.args.errs)
		})
	}
}

func TestStatServer_cleanup(t *testing.T) {
	tests := []struct {
		name string
		ss   *StatServer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ss.cleanup()
		})
	}
}

func TestStatServer_statsHandler(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	w := httptest.NewRecorder()
	tests := []struct {
		name   string
		ss     *StatServer
		args   args
		expect string
	}{
		// TODO: Add test cases.
		{
			name: "OK",
			ss: &StatServer{
				stats: &collector.CollectReport{
					Repos:    nil,
					Commits:  0,
					Lines:    0,
					Projects: 0,
				},
			},
			expect: `{"commits":0,"lines":0,"projects":0}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(tt.ss.statsHandler)
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			handler.ServeHTTP(w, req)
			if status := w.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}

			// Check the response body is what we expect.
			if eq := AreEqualJSON(w.Body.String(), tt.expect); !eq {
				t.Errorf("handler returned unexpected body: got %v want %v",
					w.Body.String(), tt.expect)
			}

		})
	}
}

func AreEqualJSON(s1, s2 string) bool {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false
	}

	return reflect.DeepEqual(o1, o2)
}
