package cache

import (
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/thales-e-security/contribstats/pkg/config"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-billy.v4/util"
	"gopkg.in/src-d/go-git-fixtures.v3"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	config.InitConfig("")
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
		{
			name: "ok",
			args: args{},
			wantGc: &GitCache{
				basepath: DefaultCache,
			},
		},
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

func TestGitCache_Add(t *testing.T) {
	td, _ := ioutil.TempDir("", "")
	td2, _ := ioutil.TempDir("", "")
	defer func() {
		os.RemoveAll(td)
		//os.RemoveAll(td2)
	}()
	type args struct {
		repo string
		url  string
	}
	tests := []struct {
		name      string
		gc        *GitCache
		args      args
		wantErr   bool
		badClone  bool
		badConfig bool
	}{
		{
			name: "OK Default",
			gc:   NewGitCache(td),
			args: args{
				repo: "github.com/unorepo/uno",
				url:  "https://github.com/unorepo/uno.git",
			},
			wantErr: false,
		},
		{
			name: "OK Empty",
			gc:   NewGitCache(""),
			args: args{
				repo: "github.com/unorepo/uno",
				url:  "https://github.com/unorepo/uno.git",
			},
			wantErr: false,
		}, {
			name: "Error BadPath",
			gc: &GitCache{
				basepath: "/jfkldasjfkladsjfkl;adsjkl;fjdsakl;jfdakls;",
			},
			args: args{
				repo: "github.com/unorepo/uno",
				url:  "https://github.com/unorepo/uno.git",
			},
			wantErr: true,
		}, {
			name: "Error BadClonin",
			gc: &GitCache{
				basepath: td2,
			},
			args: args{
				repo: "github.com/unorepo/uno",
				url:  "https://github.com/unorepo/uno.git",
			},
			badClone: true,
			wantErr:  true,
		},
		// TODO: Force a Bad Fetch to happen
		//
		//{
		//	name: "Error BadFetch",
		//	gc: &GitCache{
		//		basepath: td2,
		//	},
		//	args: args{
		//		repo: "github.com/unorepo/uno",
		//		url:  "https://github.com/unorepo/uno.git",
		//	},
		//	badConfig: true,
		//	wantErr:   true,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitDir := filepath.Join(tt.gc.Path(), tt.args.repo, ".git")
			if tt.badClone {
				// clone first, then mess it up
				tt.gc.Add(tt.args.repo, tt.args.url)
				os.RemoveAll(gitDir)
			}
			if tt.badConfig {
				tt.gc.Add(tt.args.repo, tt.args.url)
				os.Remove(filepath.Join(gitDir, "config"))
			}
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
			wantLines:   2,
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
		commit CommitIface
	}
	tests := []struct {
		name      string
		args      args
		wantLines int64
		wantErr   bool
	}{
		{
			name: "OK",
			args: args{
				commit: getGoodCommit(t),
			},
			wantLines: 1,
			wantErr:   false,
		}, {
			name: "ErrorTree",
			args: args{
				commit: &MockCommit{
					treeErr: true,
				},
			},
			wantErr: true,
		}, {
			name: "ErrorParent",
			args: args{
				commit: &MockCommit{
					treeErr:    false,
					parentsNum: 1,
					parentsErr: true,
					gc:         getGoodCommit(t),
				},
			},
			wantErr: true,
		}, {
			name: "ErrorParent",
			args: args{
				commit: &MockCommit{
					treeErr:    true,
					parentsNum: 2,
					gc:         getGoodCommit(t),
				},
			},
			wantErr: true,
		},
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

func getGoodCommit(t *testing.T) (c *object.Commit) {
	var repo *git.Repository
	var wt *git.Worktree
	var hash plumbing.Hash
	var err error
	// Clone an existing fixture repo to manipulate
	repo, err = git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL: fixtures.Basic().One().URL,
	})
	if err != nil {
		t.Error(err)
	}
	// Get the worktree to manipulate
	wt, err = repo.Worktree()
	if err != nil {
		t.Error(err)
	}

	// Create a temp file to add to it.
	err = util.WriteFile(wt.Filesystem, "foo", []byte("test"), 0755)
	if err != nil {
		t.Error(err)
	}

	//Add  path to the worktree
	_, err = wt.Add("foo")
	if err != nil {
		t.Error(err)
	}

	hash, err = wt.Commit("test", &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  "John Candy",
			Email: "john@candy.com",
			When:  time.Now(),
		},
		Committer: &object.Signature{
			Name:  "John Candy",
			Email: "john@candy.com",
			When:  time.Now(),
		},
		Parents: []plumbing.Hash{},
	})
	if err != nil {
		t.Error(err)
	}
	//hash := plumbing.NewHash("b8e471f58bcbca63b07bda20e428190409c2db47")
	c, err = repo.CommitObject(hash)
	if err != nil {
		t.Error(err)
	}
	return
}
