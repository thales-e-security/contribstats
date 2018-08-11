package server

import (
	"net/http"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"github.com/thales-e-security/contribstats/pkg/collector"
	"github.com/thales-e-security/contribstats/pkg/cache"
	"time"
	"github.com/spf13/viper"
	"encoding/json"
)

type StatServer struct {
	stats *collector.CollectReport
}

var osExit = os.Exit
var quit = make(chan os.Signal)

func NewStatServer(debug bool) (ss *StatServer) {
	ss = &StatServer{}
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// If token not provided by flag, then try by environment

	return
}

func (ss *StatServer) Start() (err error) {

	signal.Notify(quit, os.Interrupt, os.Kill)
	errs := make(chan error)
	go ss.startCollector(errs)
	// Start the Server in the background...
	go ss.startServer(errs)

	for {
		select {
		case err = <-errs:
			logrus.Error(err)
		case <-quit:
			logrus.Warn("Quitting")
			ss.cleanup()
			logrus.Exit(0)
			return
		}
	}
	return
}

func (ss *StatServer) startServer(errs chan error) {
	// Server the simple API
	http.HandleFunc("/", ss.statsHandler)
	// Start the server and wait for an error
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		errs <- err
	}
}

func (ss *StatServer) startCollector(errs chan error) () {

	var err error

	// Get a GitHubCloneCollector and make sure to add it's members and
	gh := collector.NewGitHubCloneCollector(cache.NewGitCache(cache.DefaultCache))

	// First Run....
	logrus.Info("Bootstrapping Cache and Stats")
	if ss.stats, err = gh.Collect(); err != nil {
		errs <- err
		return
	}
	logrus.Info("Updated Cache and Stats")
	// Ticker to run the job on an interval provided by the config file... defaults to 60 seconds...
	ticker := time.NewTicker(time.Duration(viper.GetInt("interval")) * time.Second)

	// Run the Collect func on a regular basis, and get ready to quit if needed
	go func() {
		for {
			select {
			case <-ticker.C:
				if ss.stats, err = gh.Collect(); err != nil {
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
	osExit(0)
}
