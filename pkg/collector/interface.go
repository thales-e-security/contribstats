//Package collector provides structs for collecting stats about contributions to git base repositories.
package collector

//Collector is a simple interface for git repo collectors for that return stats.
type Collector interface {
	// Collects stats from the API, and returns the values as a []byte of JSON content
	Collect() (stats *CollectReport, err error)
}
