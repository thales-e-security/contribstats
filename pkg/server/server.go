package server

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/thales-e-security/contribstats/pkg/cache"
	"github.com/thales-e-security/contribstats/pkg/collector"
	"github.com/thales-e-security/contribstats/pkg/config"
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
	constants config.Constants
}

var osExit = os.Exit
var cancel = make(chan struct{}, 1)
var errs = make(chan error)
var timeNewTicker = time.NewTicker

//var httpListenAndServe = http.ListenAndServe

//NewStatServer returns an instance of StatServer
func NewStatServer(constants config.Constants) (ss Server) {
	ss = &StatServer{
		stats:     nil,
		collector: collector.NewGitHubCloneCollector(constants, cache.NewGitCache(cache.DefaultCache)),
		constants: constants,
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
			return
		case <-cancel:
			ss.cleanup()
			return
		}
	}
}

func (ss *StatServer) startServer(errs chan error) {
	// Server the simple API
	mux := http.NewServeMux()
	mux.HandleFunc("/", ss.statsHandler)
	// Start the server and wait for an error
	err := http.ListenAndServe(":8080", mux)
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
	ticker := timeNewTicker(time.Duration(ss.constants.Interval) * time.Second)

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
