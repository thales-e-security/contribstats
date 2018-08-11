package collector

type Collector interface {
	// Collects stats from the API, and returns the values as a []byte of JSON content
	Collect() (stats []byte, err error)
}
