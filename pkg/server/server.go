package server

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/thales-e-security/contribstats/pkg/cache"
	"github.com/thales-e-security/contribstats/pkg/collector"
	"net/http"
	"os"
	"time"
)

// Server defines the interface for StatsServer and it's mock
type Server interface {
	Start() (err error)
}

// StatServer is starts and polls stats collection, and serves the results via a simple API
type StatServer struct {
	stats     *collector.CollectReport
	collector collector.Collector
}

var osExit = os.Exit
var cancel = make(chan struct{}, 1)
var errs = make(chan error)
var timeNewTicker = time.NewTicker
var httpListenAndServe = http.ListenAndServe

//NewStatServer returns an instance of StatServer
func NewStatServer() (ss *StatServer) {
	ss = &StatServer{
		stats:     nil,
		collector: collector.NewGitHubCloneCollector(cache.NewGitCache(cache.DefaultCache)),
	}

	return
}

//Start will start the collector and api server and then block for errors, interrupts, or cancellation
func (ss *StatServer) Start() (err error) {

	//errs := make(chan error)
	go ss.startCollector(errs)
	// Start the Server in the background...
	go ss.startServer(errs)

	for {
		select {
		case err = <-errs:
			return err
		case <-cancel:
			logrus.Warn("Got Cancel")
			ss.cleanup()
		}
	}
	return
}

func (ss *StatServer) startServer(errs chan error) {
	// Server the simple API
	mux := http.NewServeMux()
	mux.HandleFunc("/", ss.statsHandler)
	// Start the server and wait for an error
	err := httpListenAndServe(":8080", nil)
	if err != nil {
		errs <- err
	}
}

func (ss *StatServer) startCollector(errs chan error) {
	var err error
	// First Run....
	logrus.Info("Bootstrapping Cache and Stats")
	if ss.stats, err = ss.collector.Collect(); err != nil {
		errs <- err
		return
	}
	logrus.Info("Updated Cache and Stats")
	// Ticker to run the job on an interval provided by the config file... defaults to 60 seconds...
	ticker := timeNewTicker(time.Duration(viper.GetInt("interval")) * time.Second)

	// Run the Collect func on a regular basis, and get ready to quit if needed
	go func() {
		for {
			select {
			case <-ticker.C:
				if ss.stats, err = ss.collector.Collect(); err != nil {
					errs <- err
				}
				logrus.Info("Updated Cache and Stats")
			}
		}
	}()

	return
}

func (ss *StatServer) statsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ss.stats)
}

func (ss *StatServer) cleanup() {
	// TODO: Cleanup Afterwards if needed
}
