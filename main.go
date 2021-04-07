package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v32/github"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var (
	influx       string
	influxToken  string
	influxBucket string

	influxWriter api.WriteAPIBlocking
	client       *githubv4.Client
	clientv3     *github.Client
	username     string

	rootCmd = &cobra.Command{
		Use:           "gitflux",
		Short:         "Track your GitHub projects in influx",
		SilenceErrors: true,
		SilenceUsage:  true,
		// TraverseChildren:  true,
		PersistentPreRunE: initConnections,
	}
)

func initConnections(cmd *cobra.Command, args []string) error {
	var httpClient *http.Client
	token := os.Getenv("GITHUB_TOKEN")
	if len(token) == 0 {
		return fmt.Errorf("Please set your GITHUB_TOKEN env var")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	httpClient = oauth2.NewClient(context.Background(), ts)
	client = githubv4.NewClient(httpClient)

	tc := oauth2.NewClient(context.Background(), ts)
	clientv3 = github.NewClient(tc)

	var err error
	username, err = getUsername()
	if err != nil {
		return fmt.Errorf("Can't retrieve GitHub profile: %s", err)
	}

	// Create a new client using an InfluxDB server base URL and an authentication token
	idb := influxdb2.NewClient(influx, influxToken)
	// defer idb.Close()

	// Use blocking write client for writes to desired bucket
	influxWriter = idb.WriteAPIBlocking("", influxBucket)

	return nil
}

func queryWithRetry(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	if err := client.Query(context.Background(), q, variables); err != nil {
		if strings.Contains(err.Error(), "abuse-rate-limits") {
			time.Sleep(time.Minute)
			return queryWithRetry(ctx, q, variables)
		}
	}

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&influx, "influx", "http://localhost:8086", "InfluxDB address")
	rootCmd.PersistentFlags().StringVar(&influxToken, "influx-token", "", "InfluxDB auth token")
	rootCmd.PersistentFlags().StringVar(&influxBucket, "influx-bucket", "github", "InfluxDB bucket")
}
