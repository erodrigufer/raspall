package server

import (
	"fmt"
	"net/url"
	"strconv"
)

type queryOptions struct {
	removePaywall bool
	limit         int
}

func getQueryOptions(rawQuery string) (queryOptions, error) {
	options := queryOptions{}

	values, err := url.ParseQuery(rawQuery)
	if err != nil {
		return queryOptions{}, fmt.Errorf("error parsing URL query parameters: %w", err)
	}
	paywallStr := values.Get("removePaywall")
	// Default is to remove paywalled articles.
	removePaywall := true
	if paywallStr == "false" {
		removePaywall = false
	}
	options.removePaywall = removePaywall

	limitStr := values.Get("limit")
	// Default value.
	limit := 0
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			// If unable to parse value into int, use default value.
			limit = 0
		}
	}
	options.limit = limit

	return options, nil
}
