package cache

import (
	"os"
	"path/filepath"
	"strings"

	"bytes"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/format/diff"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

var DefaultCache = filepath.Join(os.TempDir(), ".ghstats", "cache")

type GitCache struct {
	basepath string
	members  []string
	domains  []string
}

func NewGitCache(basepath string) (gc *GitCache) {
	gc = &GitCache{
		basepath: basepath,
	}
	return gc
}
func (gc *GitCache) Path() string {
	return gc.basepath
}

func (gc *GitCache) Add(repo, url string) (err error) {
	repoPath := filepath.Join(gc.Path(), repo)
	var rep *git.Repository
	if _, err = os.Stat(repoPath); err != nil {
		if os.IsNotExist(err) {
			// Clone non-existing repos...
			if rep, err = git.PlainClone(repoPath, false, &git.CloneOptions{
				URL:        url,
				Progress:   &bytes.Buffer{},
				RemoteName: "",
			}); err != nil {
				return
			}
		} else {
			return
		}
		logrus.Debugf("Cloned %v into %v", repo, repoPath)
	} else {
		// Fetch existing repos to keep them up to date
		if rep, err = git.PlainOpen(repoPath); err != nil {
			return
		}
		if err = rep.Fetch(&git.FetchOptions{
			Progress: &bytes.Buffer{},
		}); err != nil {
			if git.NoErrAlreadyUpToDate.Error() == err.Error() {
				err = nil
			} else {
				return
			}
		}
		logrus.Debugf("Fetched %v in %v", repo, repoPath)
	}
	return
}

func (gc *GitCache) Stats(repo string) (commits int64, lines int64, err error) {
	logrus.Debugf("Processing repo '%s'", repo)
	var rep *git.Repository
	var logs object.CommitIter
	repoPath := filepath.Join(gc.Path(), repo)
	members := viper.GetStringSlice("members")
	domains := viper.GetStringSlice("domains")
	if rep, err = git.PlainOpen(repoPath); err != nil {
		logrus.Error(err)
		return
	}
	if logs, err = rep.Log(&git.LogOptions{
		Order: git.LogOrderDefault,
	}); err != nil {
		return
	} else {
		// For each commit entry, let's process the contents
		err = logs.ForEach(func(commit *object.Commit) (err error) {
			split := strings.Split(commit.Committer.Email, "@")
			var domain string
			if len(split) == 2 {
				domain = split[1]
			}
			// See if this commit was commited by an email address or domain we are looking for
			if stringInSlice(commit.Committer.Email, members) || stringInSlice(domain, domains) {
				// increment the commit count
				commits = commits + 1
				var newLines int64
				newLines, err = getLines(commit)
				if err != nil {
					logrus.Errorf("[%s] failed to get lines for commit %s: %s", repo, commit.Hash, err)
					return
				}
				lines = lines + newLines
				//logrus.Debugf("[%s] commit %s had %d lines", repo, commit.Hash, lines)
			}
			return
		})
		if err != nil {
			return
		}
	}
	return
}

func getLines(commit *object.Commit) (lines int64, err error) {
	// Get the lines from this commit and it's parent
	var tree *object.Tree
	var treeDiff object.Changes
	var parent *object.Commit
	var parentTree *object.Tree
	var patch *object.Patch

	// Get the tree for this commit
	if tree, err = commit.Tree(); err != nil {
		return
	}
	if commit.NumParents() != 0 {
		// Get the Parent for the commit
		if parent, err = commit.Parent(0); err != nil {
			return
		}
		// Get the Tree for the Parent commit
		if parentTree, err = parent.Tree(); err != nil {
			return
		}
	}
	// Get the Diff of the commit tree vs the parent
	if treeDiff, err = tree.Diff(parentTree); err != nil {
		return
	}
	// Get the patch of the treeDiff for processing
	if patch, err = treeDiff.Patch(); err != nil {
		return
	}
	// Iterate over the FilePatches in this diff
	for _, p := range patch.FilePatches() {
		// If it's binary in nature, let's skip it... we only want source code lines.
		if p.IsBinary() {
			continue
		}
		// Range over the chunks in a given filepatch
		for _, chunk := range p.Chunks() {
			if chunk.Type() != diff.Add {
				continue
			}
			ll := strings.Split(chunk.Content(), "\n")
			lines = lines + int64(len(ll))
		}
	}
	return
}
