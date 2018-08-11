# contribstats

Microservice/command to get a collection stats for a given GitHub Organization.

[![CircleCI](https://circleci.com/gh/thales-e-security/contribstats/tree/master.svg?style=svg)](https://circleci.com/gh/thales-e-security/contribstats/tree/master) [![codecov](https://codecov.io/gh/thales-e-security/contribstats/branch/master/graph/badge.svg)](https://codecov.io/gh/thales-e-security/contribstats) [![Go Report Card](https://goreportcard.com/badge/github.com/thales-e-security/contribstats)](https://goreportcard.com/report/github.com/thales-e-security/contribstats)


## Output

Output of the stats, will be a simple JSON format.

Example:
`
{"projects":"100","commits":"1000", "lines":"10000"}` 

## Stats

Currently the stats collected will be:

- Total \# of Projects (both contributed to and owned)
- Total \# of Commits 
- Total \# of Lines Contributed