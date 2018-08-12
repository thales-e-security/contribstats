package cache

import (
	"context"
	"github.com/pkg/errors"
	"golang.org/x/crypto/openpgp"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type MockCommit struct {
	*object.Commit
	treeErr    bool
	parentsNum int
	parentsErr bool
	gc         *object.Commit
}

func (mc *MockCommit) Tree() (*object.Tree, error) {
	if mc.treeErr {
		return nil, errors.New("Expected Error on Tree")
	}
	return mc.gc.Tree()
}

func (mc *MockCommit) PatchContext(ctx context.Context, to *object.Commit) (*object.Patch, error) {
	panic("implement me")
}

func (mc *MockCommit) Patch(to *object.Commit) (*object.Patch, error) {
	panic("implement me")
}

func (mc *MockCommit) Parents() object.CommitIter {
	panic("implement me")
}

func (mc *MockCommit) NumParents() int {
	return mc.parentsNum
}

func (mc *MockCommit) Parent(i int) (*object.Commit, error) {
	//panic("implement me")
	if mc.parentsErr {
		return nil, errors.New("Expected Error")
	}
	return mc.Commit.Parent(i)
}

func (mc *MockCommit) File(path string) (*object.File, error) {
	panic("implement me")
}

func (mc *MockCommit) Files() (*object.FileIter, error) {
	panic("implement me")
}

func (mc *MockCommit) ID() plumbing.Hash {
	panic("implement me")
}

func (mc *MockCommit) Type() plumbing.ObjectType {
	panic("implement me")
}

func (mc *MockCommit) Decode(o plumbing.EncodedObject) (err error) {
	panic("implement me")
}

func (mc *MockCommit) Encode(o plumbing.EncodedObject) error {
	panic("implement me")
}

func (mc *MockCommit) Stats() (object.FileStats, error) {
	panic("implement me")
}

func (mc *MockCommit) String() string {
	panic("implement me")
}

func (mc *MockCommit) Verify(armoredKeyRing string) (*openpgp.Entity, error) {
	panic("implement me")
}
