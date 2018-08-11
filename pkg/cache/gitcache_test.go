package cache

import (
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

var testDomains = []string{"thalesesecurity.com", "thalesesec.net", "thales-e-security.com"}
var testEmails = []string{"test@example.com", "test@example.com"}

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	viper.Set("members", []string{"scates@thalesesec.net"})
	viper.Set("domains", []string{"thalesesec.net"})
}

func TestNewGitCache(t *testing.T) {
	type args struct {
		basepath string
	}
	tests := []struct {
		name   string
		args   args
		wantGc *GitCache
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotGc := NewGitCache(tt.args.basepath); !reflect.DeepEqual(gotGc, tt.wantGc) {
				t.Errorf("NewGitCache() = %v, want %v", gotGc, tt.wantGc)
			}
		})
	}
}

func TestGitCache_Path(t *testing.T) {
	tests := []struct {
		name string
		gc   *GitCache
		want string
	}{
		{
			name: "OK",
			gc:   NewGitCache(DefaultCache),
			want: DefaultCache,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.gc.Path(); got != tt.want {
				t.Errorf("GitCache.Path() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitCache_Repos(t *testing.T) {
	tests := []struct {
		name      string
		gc        *GitCache
		wantRepos []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRepos := tt.gc.Repos(); !reflect.DeepEqual(gotRepos, tt.wantRepos) {
				t.Errorf("GitCache.Repos() = %v, want %v", gotRepos, tt.wantRepos)
			}
		})
	}
}

func TestGitCache_Add(t *testing.T) {
	type args struct {
		repo string
		url  string
	}
	tests := []struct {
		name    string
		gc      *GitCache
		args    args
		wantErr bool
	}{
		{
			name: "OK",
			gc:   NewGitCache(DefaultCache),
			args: args{
				repo: "github.com/unorepo/uno",
				url:  "https://github.com/unorepo/uno.git",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.gc.Add(tt.args.repo, tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("GitCache.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGitCache_Stats(t *testing.T) {
	type args struct {
		repo    string
		members []string
	}
	tests := []struct {
		name        string
		gc          *GitCache
		args        args
		wantCommits int64
		wantLines   int64
		wantErr     bool
	}{
		{
			name: "OK",
			gc:   NewGitCache(DefaultCache),
			args: args{
				repo: "github.com/unorepo/uno",
			},
			wantCommits: 1,
			wantLines:   0,
			wantErr:     false,
		}, {
			name: "Error",
			gc:   NewGitCache(DefaultCache),
			args: args{
				repo: "github.com/notreallyhere/repo",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCommits, gotLines, err := tt.gc.Stats(tt.args.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("GitCache.Stats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCommits != tt.wantCommits {
				t.Errorf("GitCache.Stats() gotCommits = %v, want %v", gotCommits, tt.wantCommits)
			}
			if gotLines != tt.wantLines {
				t.Errorf("GitCache.Stats() gotLines = %v, want %v", gotLines, tt.wantLines)
			}
		})
	}
}

func Test_getLines(t *testing.T) {
	type args struct {
		commit *object.Commit
	}

	tests := []struct {
		name      string
		args      args
		wantLines int64
		wantErr   bool
	}{
		// TODO: Create cases
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLines, err := getLines(tt.args.commit)
			if (err != nil) != tt.wantErr {
				t.Errorf("getLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotLines != tt.wantLines {
				t.Errorf("getLines() = %v, want %v", gotLines, tt.wantLines)
			}
		})
	}
}
