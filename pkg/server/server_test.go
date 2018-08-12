package server

import (
	"reflect"
	"testing"

	"encoding/json"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/thales-e-security/contribstats/pkg/cache"
	"github.com/thales-e-security/contribstats/pkg/collector"
	"net/http"
	"net/http/httptest"
	"time"
)

func init() {
	viper.Set("interval", 60)
}

func TestNewStatServer(t *testing.T) {

	tests := []struct {
		name   string
		wantSs *StatServer
	}{
		{
			name: "OK",
			wantSs: &StatServer{
				collector: collector.NewGitHubCloneCollector(cache.NewGitCache(cache.DefaultCache)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSs := NewStatServer(); !reflect.DeepEqual(gotSs, tt.wantSs) {
				t.Errorf("NewStatServer() = %v, want %v", gotSs, tt.wantSs)
			}
		})
	}
}

func TestStatServer_Start(t *testing.T) {
	tests := []struct {
		name    string
		ss      *StatServer
		cancel  bool
		quit    bool
		error   bool
		wantErr bool
	}{
		{
			name:    "OK",
			ss:      NewStatServer(),
			wantErr: false,
		},
		{
			name:    "Error",
			ss:      NewStatServer(),
			wantErr: true,
		},
		{
			name:   "Cancel",
			ss:     NewStatServer(),
			cancel: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotExitCode int
			// store old os.Exit
			oldOSExit := osExit
			// override os.Exit
			osExit = func(code int) {
				gotExitCode = code
			}
			if tt.wantErr {
				httpListenAndServe = func(addr string, handler http.Handler) (err error) {
					err = errors.New("expected error")
					return
				}
			}

			// Start the server
			go tt.ss.Start()
			// wait for it...
			time.Sleep(10 * time.Millisecond)

			// Canceling
			if tt.cancel {
				// Kill it ...
				cancel <- struct{}{}
			}

			// repair os.Exit
			osExit = oldOSExit

			// See what we gotExitCode
			if gotExitCode != 0 {
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
		//TODO: Add Test Cases
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
		name    string
		ss      *StatServer
		args    args
		wantErr bool
	}{
		{
			name: "OK",
			ss:   NewStatServer(),
			args: args{
				errs: errs,
			},
			wantErr: false,
		},
		{
			name: "Error",
			ss: &StatServer{
				collector: &MockCollector{
					wantErr: true,
				},
			},
			args: args{
				errs: errs,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := make(chan time.Time)

			timeNewTicker = func(d time.Duration) *time.Ticker {
				return &time.Ticker{
					C: c,
				}
			}

			go tt.ss.startCollector(tt.args.errs)
			go func() {
				time.Sleep(10 * time.Millisecond)
				c <- time.Now()
			}()
			select {
			case <-c:

			}
		})
	}
}

func TestStatServer_cleanup(t *testing.T) {
	tests := []struct {
		name string
		ss   *StatServer
	}{
		{
			name: "OK",
			ss:   &StatServer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			osExit = func(code int) {

			}
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

type MockCollector struct {
	wantErr bool
}

func (mc *MockCollector) Collect() (stats *collector.CollectReport, err error) {
	stats = &collector.CollectReport{}
	if mc.wantErr {
		err = errors.New("expected error")
	}
	return
}
