package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/spf13/cobra"
)

var (
	repositoryCmd = &cobra.Command{
		Use:   "repository [foo/bar]",
		Short: "tracks repositories",
		RunE: func(cmd *cobra.Command, args []string) error {
			var repos []Repo

			if len(args) > 0 {
				if strings.Contains(args[0], "/") {
					as := strings.Split(args[0], "/")
					r, err := repository(as[0], as[1])
					if err != nil {
						return err
					}

					repos = append(repos, r)
				} else {
					return fmt.Errorf("invalid repository: %s", args[0])
				}
			} else {
				// fetch all user source repositories per default
				fmt.Printf("Finding user's source repos...\n")
				var err error
				repos, err = repositories()
				if err != nil {
					return err
				}
				fmt.Printf("Found %d repos\n", len(repos))
			}

			return parseRepos(repos, influxWriter)
		},
	}
)

func parseRepos(repos []Repo, w api.WriteAPIBlocking) error {
	for _, r := range repos {
		fmt.Printf("Parsing %s\n", r.NameWithOwner)

		p := influxdb2.NewPointWithMeasurement("stars").
			AddTag("repo", "github.com/"+r.NameWithOwner).
			AddField("value", r.Stargazers).
			SetTime(time.Now())
		err := w.WritePoint(context.Background(), p)
		if err != nil {
			return err
		}

		p = influxdb2.NewPointWithMeasurement("watchers").
			AddTag("repo", "github.com/"+r.NameWithOwner).
			AddField("value", r.Watchers).
			SetTime(time.Now())
		err = w.WritePoint(context.Background(), p)
		if err != nil {
			return err
		}

		p = influxdb2.NewPointWithMeasurement("forks").
			AddTag("repo", "github.com/"+r.NameWithOwner).
			AddField("value", r.Forks).
			SetTime(time.Now())
		err = w.WritePoint(context.Background(), p)
		if err != nil {
			return err
		}

		p = influxdb2.NewPointWithMeasurement("commits").
			AddTag("repo", "github.com/"+r.NameWithOwner).
			AddField("value", r.Commits).
			SetTime(time.Now())
		err = w.WritePoint(context.Background(), p)
		if err != nil {
			return err
		}

		// parse PRs
		err = parsePullRequests(r, w)
		if err != nil {
			return err
		}
		// parse issues
		err = parseIssues(r, w)
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(repositoryCmd)
}
