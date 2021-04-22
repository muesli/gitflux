package main

import (
	"context"

	"github.com/shurcooL/githubv4"
)

var sponsorsQuery struct {
	Viewer struct {
		SponsorshipAsMaintainer struct {
			TotalCount githubv4.Int
		} `graphql:"sponsorshipsAsMaintainer(first: 100, orderBy: {field: CREATED_AT, direction: DESC})"`
	}
}

func sponsors() (int, error) {
	if err := queryWithRetry(context.Background(), &sponsorsQuery, map[string]interface{}{}); err != nil {
		return 0, err
	}

	return int(sponsorsQuery.Viewer.SponsorshipAsMaintainer.TotalCount), nil
}
