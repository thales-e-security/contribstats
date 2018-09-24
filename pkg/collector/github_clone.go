package collector

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/thales-e-security/contribstats/pkg/cache"
	"github.com/thales-e-security/contribstats/pkg/config"
)

var timeAfter = time.After

func init() {
	//logrus.SetLevel(logrus.DebugLevel)
}

//GitHubCloneCollector uses offline caching of repos to obtain stats of contribution by members and domains of interest
type GitHubCloneCollector struct {
	client    *github.Client
	cache     cache.Cache
	ctx       context.Context
	done      chan *RepoResults
	errs      chan error
	constants config.Config
}

//NewGitHubCloneCollector returns a GitHubCloneCollector.
//
func NewGitHubCloneCollector(contants config.Config, c cache.Cache) (ghc *GitHubCloneCollector) {
	ghc = &GitHubCloneCollector{
		cache:     c,
		constants: contants,
	}
	// Set the Client
	ghc.client, ghc.ctx = NewV3Client(contants)
	return
}

//RepoResults contains results from an individual repository
type RepoResults struct {
	Repo    string `json:"repo"`
	Commits int64  `json:"commits"`
	Lines   int64  `json:"lines"`
}

//CollectReport contains the results of an entire collection of repos, and an aggregated value of each stats
type CollectReport struct {
	Repos    []*RepoResults `json:"repos,omitempty"`
	Commits  int64          `json:"commits"`
	Lines    int64          `json:"lines"`
	Projects int64          `json:"projects"`
}

//Collect iterates over all members in the organization to aggregate their OpenSource contributions offline
func (ghc *GitHubCloneCollector) Collect() (stats *CollectReport, err error) {
	var repos []*github.Repository
	var done = make(chan *RepoResults)
	var errs = make(chan error)
	stats = &CollectReport{}
	// List of Repos
	for _, org := range ghc.constants.Organizations {
		var tRepos []*github.Repository
		if tRepos, _, err = ghc.client.Repositories.ListByOrg(ghc.ctx, org, nil); err != nil {
			return
		}
		repos = append(repos, tRepos...)
	}
	go func() {
		for _, repo := range repos {
			go ghc.processRepo(repo, done, errs)
		}
	}()
	// Drain the channels
	for i := 1; i <= len(repos); i++ {
		select {
		case d := <-done:
			logrus.Debugf("Done with: %v", d.Repo)
			stats.Repos = append(stats.Repos, d)
			stats.Commits = stats.Commits + d.Commits
			stats.Lines = stats.Lines + d.Lines
		case err := <-errs:
			logrus.Error(err)
		case <-timeAfter(10 * time.Minute):
			return nil, errors.New("Timed out")
		}
	}
	// For convenience, return a count of repos
	stats.Projects = int64(len(stats.Repos))
	logrus.Debugf("Finished Collecting Stats")
	// TODO: return and/or display the stats
	return
}

//TODO: Process activity on a given repo for stats from this organization.
func (ghc *GitHubCloneCollector) processRepo(repo *github.Repository, done chan *RepoResults, errs chan error) {

	// Check if repo is on blacklist
	for _, reponame := range ghc.constants.Blacklist {
		if strings.ToLower(reponame) == strings.ToLower(repo.String()) {
			logrus.Warn("Skipping Blacklisted item: ")
		}
	}

	var err error
	// First let's clone it to the local cache dir.
	name := filepath.Join("github.com", repo.GetFullName())
	if err = ghc.cache.Add(name, repo.GetCloneURL()); err != nil {
		err = errors.Wrap(err, "add")
		errs <- err
		return
	}

	// Generate the scaffolding
	r := &RepoResults{
		Repo: name,
	}
	// Get Stats on cached repo...
	if r.Commits, r.Lines, err = ghc.cache.Stats(name); err != nil {
		err = errors.Wrap(err, "stats")
		errs <- err
		return
	}
	done <- r
}
