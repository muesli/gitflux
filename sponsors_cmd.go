package main

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/spf13/cobra"
)

var sponsorsCmd = &cobra.Command{
	Use:   "sponsors",
	Short: "tracks user sponsors",
	RunE: func(cmd *cobra.Command, args []string) error {
		return parseSponsors(influxWriter)
	},
}

func parseSponsors(w api.WriteAPIBlocking) error {
	fmt.Printf("Finding sponsors for user...\n")
	s, err := sponsors()
	if err != nil {
		return err
	}
	fmt.Printf("Found %d sponsors\n", s)

	p := influxdb2.NewPointWithMeasurement("sponsors").
		AddTag("user", username).
		AddField("value", s).
		SetTime(time.Now())
	return w.WritePoint(context.Background(), p)
}

func init() {
	rootCmd.AddCommand(sponsorsCmd)
}
