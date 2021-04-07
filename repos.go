package main

import (
	"context"
	"time"

	"github.com/shurcooL/githubv4"
)

var reposQuery struct {
	User struct {
		Login        githubv4.String
		Repositories struct {
			TotalCount githubv4.Int
			Edges      []struct {
				Cursor githubv4.String
				Node   QLRepository
			}
		} `graphql:"repositories(first: 100, after:$after privacy: PUBLIC, isFork: false, ownerAffiliations: OWNER, orderBy: {field: CREATED_AT, direction: DESC})"`
	} `graphql:"repositoryOwner(login:$username)"`
}

var repoQuery struct {
	Repository QLRepository `graphql:"repository(owner: $owner, name: $name)"`
}

func repository(owner string, name string) (Repo, error) {
	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(name),
	}

	if err := queryWithRetry(context.Background(), &repoQuery, variables); err != nil {
		return Repo{}, err
	}

	return RepoFromQL(repoQuery.Repository), nil
}

func repositories(owner string) ([]Repo, error) {
	var after *githubv4.String
	var repos []Repo

	for {
		variables := map[string]interface{}{
			"username": githubv4.String(owner),
			"after":    after,
		}

		if err := queryWithRetry(context.Background(), &reposQuery, variables); err != nil {
			return nil, err
		}
		if len(reposQuery.User.Repositories.Edges) == 0 {
			break
		}

		for _, v := range reposQuery.User.Repositories.Edges {
			repos = append(repos, RepoFromQL(v.Node))

			after = &v.Cursor
		}
	}

	return repos, nil
}

type QLRepository struct {
	Owner struct {
		Login githubv4.String
	}
	Name           githubv4.String
	NameWithOwner  githubv4.String
	URL            githubv4.String
	Description    githubv4.String
	IsPrivate      githubv4.Boolean
	ForkCount      githubv4.Int
	StargazerCount githubv4.Int

	Watchers struct {
		TotalCount githubv4.Int
	}

	BranchEntity struct {
		Commits struct {
			History struct {
				TotalCount githubv4.Int
			}
		} `graphql:"... on Commit"`
	} `graphql:"object(expression: \"HEAD\")"`
}

type Repo struct {
	Owner         string
	Name          string
	NameWithOwner string
	URL           string
	Description   string
	Stargazers    int
	Watchers      int
	Forks         int
	Commits       int
	LastRelease   Release
}

type Release struct {
	Name        string
	TagName     string
	PublishedAt time.Time
	URL         string
}

func RepoFromQL(repo QLRepository) Repo {
	return Repo{
		Owner:         string(repo.Owner.Login),
		Name:          string(repo.Name),
		NameWithOwner: string(repo.NameWithOwner),
		URL:           string(repo.URL),
		Description:   string(repo.Description),
		Stargazers:    int(repo.StargazerCount),
		Watchers:      int(repo.Watchers.TotalCount),
		Forks:         int(repo.ForkCount),
		Commits:       int(repo.BranchEntity.Commits.History.TotalCount),
	}
}
