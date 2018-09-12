package server

import (
	"context"
	"encoding/json"
	"github.com/rs/cors"
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
var cancel = make(chan bool, 1)
var collectorCancel = make(chan bool)
var serverCancel = make(chan bool)
var errs = make(chan error)
var timeNewTicker = time.NewTicker

//var httpListenAndServe = http.ListenAndServe

//NewStatServer returns an instance of StatServer
func NewStatServer(constants config.Constants) (ss Server) {
	if constants.Cache == "" {
		constants.Cache = cache.DefaultCache
	}
	ss = &StatServer{
		collector: collector.NewGitHubCloneCollector(constants, cache.NewGitCache(constants.Cache)),
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
			serverCancel <- true
			collectorCancel <- true
			return
		}
	}
}

func (ss *StatServer) startServer(errs chan error) {
	// Server the simple API
	mux := http.NewServeMux()
	mux.HandleFunc("/", ss.statsHandler)
	ctx := context.Background()
	// Handler
	var handler http.Handler
	var c *cors.Cors
	if ss.constants.Origins != nil {
		c = cors.New(cors.Options{
			AllowedOrigins: ss.constants.Origins,
		})
	} else {
		c = cors.Default()
	}
	handler = c.Handler(mux)
	// Start the server and wait for an error
	svr := http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	go func() {
		err := svr.ListenAndServe()
		if err != nil {
			errs <- err
		}
	}()

	for {
		select {
		case <-errs:
			return
		case <-serverCancel:
			svr.Shutdown(ctx)
			return
		}
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
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(ss.stats)
}

func (ss *StatServer) cleanup() {
	// TODO: Cleanup Afterwards if needed
}
