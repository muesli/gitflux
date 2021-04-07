package main

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/shurcooL/githubv4"
)

type QLLabel struct {
	Name githubv4.String
}

type QLIssue struct {
	Title     githubv4.String
	State     githubv4.IssueState
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

type Issue struct {
	Title     string
	Open      bool
	Assignees []User
	Labels    []string
}

var issuesQuery struct {
	Repository struct {
		Issues struct {
			TotalCount githubv4.Int
			Edges      []struct {
				Cursor githubv4.String
				Node   QLIssue
			}
		} `graphql:"issues(first: 100, after:$after)"`

		/*
			OpenIssues struct {
				TotalCount githubv4.Int
			} `graphql:"filterBy: {states: OPEN}"`
		*/
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func issues(owner string, name string) ([]Issue, error) {
	var after *githubv4.String
	var issues []Issue

	for {
		variables := map[string]interface{}{
			"owner": githubv4.String(owner),
			"name":  githubv4.String(name),
			"after": after,
		}

		if err := queryWithRetry(context.Background(), &issuesQuery, variables); err != nil {
			return nil, err
		}
		if len(issuesQuery.Repository.Issues.Edges) == 0 {
			break
		}

		for _, v := range issuesQuery.Repository.Issues.Edges {
			issues = append(issues, IssueFromQL(v.Node))

			after = &v.Cursor
		}
	}

	return issues, nil
}

func IssueFromQL(issue QLIssue) Issue {
	var labels []string
	for _, v := range issue.Labels.Edges {
		labels = append(labels, string(v.Node.Name))
	}

	var assignees []User
	for _, v := range issue.Assignees.Edges {
		assignees = append(assignees, UserFromQL(v.Node))
	}

	return Issue{
		Title:     string(issue.Title),
		Open:      issue.State == githubv4.IssueStateOpen,
		Assignees: assignees,
		Labels:    labels,
	}
}

func parseIssues(r Repo, w api.WriteAPIBlocking) error {
	fmt.Printf("\tFinding issues for repo...\n")
	is, err := issues(r.Owner, r.Name)
	if err != nil {
		return err
	}
	fmt.Printf("\tFound %d issues!\n", len(is))

	open := 0
	assignees := map[string]int{}
	assigneesOpen := map[string]int{}
	labels := map[string]int{}
	labelsOpen := map[string]int{}
	for _, v := range is {
		if v.Open {
			open++
			for _, l := range v.Labels {
				labelsOpen[l]++
			}
			for _, a := range v.Assignees {
				assigneesOpen[a.Login]++
			}
		}

		for _, l := range v.Labels {
			labels[l]++
		}
		for _, a := range v.Assignees {
			assignees[a.Login]++
		}
	}

	p := influxdb2.NewPointWithMeasurement("issues").
		AddTag("repo", "github.com/"+r.NameWithOwner).
		AddField("total", len(is)).
		AddField("open", open).
		SetTime(time.Now())
	err = w.WritePoint(context.Background(), p)
	if err != nil {
		return err
	}

	for k, v := range labels {
		p = influxdb2.NewPointWithMeasurement("issues_labels").
			AddTag("repo", "github.com/"+r.NameWithOwner).
			AddTag("name", k).
			AddField("total", v).
			AddField("open", labelsOpen[k]).
			SetTime(time.Now())
		err = w.WritePoint(context.Background(), p)
		if err != nil {
			return err
		}
	}

	for k, v := range assignees {
		p = influxdb2.NewPointWithMeasurement("issues_assignees").
			AddTag("repo", "github.com/"+r.NameWithOwner).
			AddTag("name", k).
			AddField("total", v).
			AddField("open", assigneesOpen[k]).
			SetTime(time.Now())
		err = w.WritePoint(context.Background(), p)
		if err != nil {
			return err
		}
	}

	return nil
}
