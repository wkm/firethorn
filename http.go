package firethorn

import (
	"strings"
)

// parses a query
func ParseInsert(query string) []Query {
	queryStrings := strings.Split(query, "---")
	queries := make([]Query, len(queryStrings))

	for _, query := range queryStrings {
		print("parsing query: " + query)
	}

	return queries
}
