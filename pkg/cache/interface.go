package cache

type Cache interface {
	Path() string
	Add(repo, url string) (err error)
	Stats(repo string) (commits int64, lines int64, err error)
}
