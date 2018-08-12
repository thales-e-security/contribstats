package collector

import (
	"reflect"
	"testing"

	"fmt"
	"github.com/google/go-github/github"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/thales-e-security/contribstats/pkg/cache"
	"os"
	"time"
)

var testDomains = []string{"thalesesecurity.com", "thalesesec.net", "thales-e-security.com"}
var testEmails = []string{"test@example.com", "test@example.com"}
var testCache = cache.NewGitCache(cache.DefaultCache)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	//
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	viper.SetEnvPrefix("CONTRIBSTATS")
	viper.AutomaticEnv() // read in environment variables that match
	viper.AddConfigPath(home)
	viper.AddConfigPath(".")
	viper.AddConfigPath("/config")
	viper.SetConfigName(".contribstats")
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
		wantTimeout   bool
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
		}, {
			name:          "Timeout",
			ghc:           NewGitHubCloneCollector(testCache),
			wantStats:     false,
			wantErr:       true,
			wantTimeout:   true,
			organizations: []string{"thales-e-security"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.organizations != nil {
				viper.Set("organizations", tt.organizations)
			}
			if tt.wantTimeout {
				timeAfter = func(d time.Duration) <-chan time.Time {
					return time.After(time.Millisecond)
				}
			} else {
				timeAfter = time.After
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
		}, {
			name: "Error Add",
			ghc:  NewGitHubCloneCollector(&MockCache{add: true}),
			args: args{
				repo: &github.Repository{
					Name:     github.String("linux-kernel"),
					FullName: github.String("thales-e-security/linux-kernel"),
				},
				done: make(chan *RepoResults, 1),
				errs: make(chan error, 1),
			},
			wantErr: true,
		}, {
			name: "Error Add",
			ghc:  NewGitHubCloneCollector(&MockCache{stats: true}),
			args: args{
				repo: &github.Repository{
					Name:     github.String("linux-kernel"),
					FullName: github.String("thales-e-security/linux-kernel"),
				},
				done: make(chan *RepoResults, 1),
				errs: make(chan error, 1),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ghc.processRepo(tt.args.repo, tt.args.done, tt.args.errs)
			select {
			case err := <-tt.args.errs:
				if (err != nil) != tt.wantErr {
					t.Errorf("GitHubCloneCollector.processRepo() error = %v, wantErr %v", (err != nil), tt.wantErr)
					return
				}
			case <-tt.args.done:

			}

		})
	}
}

type MockCache struct {
	add   bool
	stats bool
}

func (mc *MockCache) Path() string {
	panic("implement me")
}

func (mc *MockCache) Add(repo, url string) (err error) {
	if mc.add {
		err = errors.New("expected error")

	}
	return
}

func (mc *MockCache) Stats(repo string) (commits int64, lines int64, err error) {
	if mc.stats {
		err = errors.New("expected error")

	}
	return
}
