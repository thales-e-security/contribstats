# contribstats

Microservice/command to get a collection stats for a given GitHub Organization.

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