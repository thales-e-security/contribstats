package collector

import (
	"reflect"
	"testing"

	"github.com/google/go-github/github"
	"github.com/thales-e-security/contribstats/pkg/cache"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var testDomains = []string{"thalesesecurity.com", "thalesesec.net", "thales-e-security.com"}
var testEmails = []string{"test@example.com", "test@example.com"}
var testCache = cache.NewGitCache(cache.DefaultCache)

func init() {
	logrus.SetLevel(logrus.DebugLevel)

	//
	viper.Set("organizations", []string{"thales-e-security"})
}

func TestNewGitHubCloneCollector(t *testing.T) {
	gh := NewGitHubCloneCollector(testCache)

	type args struct {
		organization string
		token        string
		cache        cache.Cache
	}
	tests := []struct {
		name string
		args args
		want *GitHubCloneCollector
	}{
		{
			name: "thales-e-security",
			args: args{
				cache: testCache,
			},
			want: gh,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGitHubCloneCollector(tt.args.cache); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGitHubCloneCollector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitHubCloneCollector_Collect(t *testing.T) {

	tests := []struct {
		name          string
		ghc           *GitHubCloneCollector
		wantStats     bool
		wantErr       bool
		organizations []string
	}{

		{
			name:          "OK",
			ghc:           NewGitHubCloneCollector(testCache),
			wantStats:     true,
			wantErr:       false,
			organizations: []string{"thales-e-security"},
		},
		{
			name:          "Error",
			ghc:           NewGitHubCloneCollector(testCache),
			wantStats:     true,
			wantErr:       true,
			organizations: []string{"tthales-e-security"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.organizations != nil {
				viper.Set("organizations", tt.organizations)
			}
			gotStats, err := tt.ghc.Collect()
			if (err != nil) != tt.wantErr {
				t.Errorf("GitHubCloneCollector.Collect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (gotStats != nil) != tt.wantStats {
				t.Errorf("GitHubCloneCollector.Collect() stats = %v, wantStats %v", (gotStats != nil), tt.wantStats)
				return
			}

		})
	}
}

func TestGitHubCloneCollector_processRepo(t *testing.T) {
	type args struct {
		repo *github.Repository
		done chan *RepoResults
		errs chan error
	}
	tests := []struct {
		name    string
		ghc     *GitHubCloneCollector
		args    args
		wantErr bool
	}{
		{
			name: "Good",
			ghc:  NewGitHubCloneCollector(cache.NewGitCache(cache.DefaultCache)),
			args: args{
				repo: &github.Repository{
					Name:     github.String("linux-kernel"),
					FullName: github.String("thales-e-security/linux-kernel"),
				},
				done: make(chan *RepoResults, 1),
				errs: make(chan error, 1),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ghc.processRepo(tt.args.repo, tt.args.done, tt.args.errs)

			select {
			case err := <-tt.args.errs:
				if err != nil {
					t.Errorf("GitHubCloneCollector.processRepo() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			case <-tt.args.done:

			}

		})
	}
}
