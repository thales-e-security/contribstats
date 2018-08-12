package collector

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/thales-e-security/contribstats/pkg/cache"
)

func init() {
	//logrus.SetLevel(logrus.DebugLevel)
}

//GitHubCloneCollector uses offline caching of repos to obtain stats of contribution by members and domains of interest
type GitHubCloneCollector struct {
	client *github.Client
	cache  cache.Cache
	ctx    context.Context
}

//NewGitHubCloneCollector returns a GitHubCloneCollector.
//
func NewGitHubCloneCollector(c cache.Cache) (ghc *GitHubCloneCollector) {
	ghc = &GitHubCloneCollector{
		cache: c,
	}
	// Set the Client
	ghc.client, ghc.ctx = NewV3Client()
	return
}

//RepoResults contains results from an individual repository
type RepoResults struct {
	Repo    string
	Commits int64
	Lines   int64
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
	for _, org := range viper.GetStringSlice("organizations") {
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
		case <-time.After(10 * time.Minute):
			// TODO: a senstible timeout for when running in parallel
			logrus.Fatal("Timeed out")
			break
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

	var err error
	// First let's clone it to the local cache dir.
	name := filepath.Join("github.com", repo.GetFullName())
	if err := ghc.cache.Add(name, repo.GetCloneURL()); err != nil {
		errors.Wrap(err, fmt.Sprintf("%v:", name))
		errs <- err
		return
	}
	// Generate the scaffolding
	r := &RepoResults{
		Repo: name,
	}
	// Get Stats on cached repo...
	r.Commits, r.Lines, err = ghc.cache.Stats(name)
	if err != nil {
		errors.Wrap(err, fmt.Sprintf("%v:", name))
		errs <- err
		return
	}
	done <- r
}
