package cache

import (
	"context"
	"golang.org/x/crypto/openpgp"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

//Cache defines how all caching backends must behave.  Caches are required to return stats.
type Cache interface {
	Path() string
	Add(repo, url string) (err error)
	Stats(repo string) (commits int64, lines int64, err error)
}

//CommitIface is interface for Commits since go-git doesn't provide an interface.
type CommitIface interface {
	// Tree returns the Tree from the commit.
	Tree() (*object.Tree, error)
	// Patch returns the Patch between the actual commit and the provided one.
	// Error will be return if context expires. Provided context must be non-nil
	PatchContext(ctx context.Context, to *object.Commit) (*object.Patch, error)
	// Patch returns the Patch between the actual commit and the provided one.
	Patch(to *object.Commit) (*object.Patch, error)
	// Parents return a CommitIter to the parent Commits.
	Parents() object.CommitIter
	// NumParents returns the number of parents in a commit.
	NumParents() int
	// Parent returns the ith parent of a commit.
	Parent(i int) (*object.Commit, error)
	// File returns the file with the specified "path" in the commit and a
	// nil error if the file exists. If the file does not exist, it returns
	// a nil file and the ErrFileNotFound error.
	File(path string) (*object.File, error)
	// Files returns a FileIter allowing to iterate over the Tree
	Files() (*object.FileIter, error)
	// ID returns the object ID of the commit. The returned value will always match
	// the current value of Commit.Hash.
	//
	// ID is present to fulfill the Object interface.
	ID() plumbing.Hash
	// Type returns the type of object. It always returns plumbing.CommitObject.
	//
	// Type is present to fulfill the Object interface.
	Type() plumbing.ObjectType
	// Decode transforms a plumbing.EncodedObject into a Commit struct.
	Decode(o plumbing.EncodedObject) (err error)
	// Encode transforms a Commit into a plumbing.EncodedObject.
	Encode(o plumbing.EncodedObject) error
	// Stats shows the status of commit.
	Stats() (object.FileStats, error)
	String() string
	// Verify performs PGP verification of the commit with a provided armored
	// keyring and returns openpgp.Entity associated with verifying key on success.
	Verify(armoredKeyRing string) (*openpgp.Entity, error)
}
