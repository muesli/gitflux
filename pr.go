package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/shurcooL/githubv4"
)

type QLPullRequest struct {
	Title     githubv4.String
	State     githubv4.PullRequestState
	Assignees struct {
		Edges []struct {
			Node QLUser
		}
	} `graphql:"assignees(first: 100)"`
	Labels struct {
		Edges []struct {
			Node QLLabel
		}
	} `graphql:"labels(first: 100)"`
}

type PullRequest struct {
	Title     string
	Open      bool
	Merged    bool
	Assignees []User
	Labels    []string
}

var pullRequestsQuery struct {
	Repository struct {
		PullRequests struct {
			TotalCount githubv4.Int
			Edges      []struct {
				Cursor githubv4.String
				Node   QLPullRequest
			}
		} `graphql:"pullRequests(first: 100, after:$after)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func pullRequests(owner string, name string) ([]PullRequest, error) {
	var after *githubv4.String
	var prs []PullRequest

	for {
		variables := map[string]interface{}{
			"owner": githubv4.String(owner),
			"name":  githubv4.String(name),
			"after": after,
		}

		err := client.Query(context.Background(), &pullRequestsQuery, variables)
		if err != nil {
			if strings.Contains(err.Error(), "abuse-rate-limits") {
				time.Sleep(time.Minute)
				continue
			}
			return nil, err
		}
		if len(pullRequestsQuery.Repository.PullRequests.Edges) == 0 {
			break
		}

		for _, v := range pullRequestsQuery.Repository.PullRequests.Edges {
			prs = append(prs, PullRequestFromQL(v.Node))

			after = &v.Cursor
		}
	}

	return prs, nil
}

func PullRequestFromQL(pr QLPullRequest) PullRequest {
	var labels []string
	for _, v := range pr.Labels.Edges {
		labels = append(labels, string(v.Node.Name))
	}

	var assignees []User
	for _, v := range pr.Assignees.Edges {
		assignees = append(assignees, UserFromQL(v.Node))
	}

	return PullRequest{
		Title:     string(pr.Title),
		Open:      pr.State == githubv4.PullRequestStateOpen,
		Merged:    pr.State == githubv4.PullRequestStateMerged,
		Assignees: assignees,
		Labels:    labels,
	}
}

func parsePullRequests(r Repo, w api.WriteAPIBlocking) error {
	fmt.Printf("\tFinding PRs for repo...\n")
	prs, err := pullRequests(r.Owner, r.Name)
	if err != nil {
		return err
	}
	fmt.Printf("\tFound %d PRs!\n", len(prs))

	open := 0
	merged := 0
	assignees := map[string]int{}
	assigneesOpen := map[string]int{}
	assigneesMerged := map[string]int{}
	labels := map[string]int{}
	labelsOpen := map[string]int{}
	labelsMerged := map[string]int{}
	for _, v := range prs {
		if v.Open {
			open++
			for _, l := range v.Labels {
				labelsOpen[l]++
			}
			for _, a := range v.Assignees {
				assigneesOpen[a.Login]++
			}
		}
		if v.Merged {
			merged++
			for _, l := range v.Labels {
				labelsMerged[l]++
			}
			for _, a := range v.Assignees {
				assigneesMerged[a.Login]++
			}
		}

		for _, l := range v.Labels {
			labels[l]++
		}
		for _, a := range v.Assignees {
			assignees[a.Login]++
		}
	}

	p := influxdb2.NewPointWithMeasurement("pullrequests").
		AddTag("repo", "github.com/"+r.NameWithOwner).
		AddField("total", len(prs)).
		AddField("open", open).
		AddField("merged", merged).
		SetTime(time.Now())
	err = w.WritePoint(context.Background(), p)
	if err != nil {
		return err
	}

	for k, v := range labels {
		p = influxdb2.NewPointWithMeasurement("pullrequests_labels").
			AddTag("repo", "github.com/"+r.NameWithOwner).
			AddTag("name", k).
			AddField("total", v).
			AddField("merged", labelsMerged[k]).
			AddField("open", labelsOpen[k]).
			SetTime(time.Now())
		err = w.WritePoint(context.Background(), p)
		if err != nil {
			return err
		}
	}

	for k, v := range assignees {
		p = influxdb2.NewPointWithMeasurement("pullrequests_assignees").
			AddTag("repo", "github.com/"+r.NameWithOwner).
			AddTag("name", k).
			AddField("total", v).
			AddField("merged", assigneesMerged[k]).
			AddField("open", assigneesOpen[k]).
			SetTime(time.Now())
		err = w.WritePoint(context.Background(), p)
		if err != nil {
			return err
		}
	}

	return nil
}
