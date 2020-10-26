package main

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/spf13/cobra"
)

var (
	relationshipsCmd = &cobra.Command{
		Use:   "relationships",
		Short: "tracks user relationships",
		RunE: func(cmd *cobra.Command, args []string) error {
			return parseRelationships(influxWriter)
		},
	}
)

func parseRelationships(w api.WriteAPIBlocking) error {
	fmt.Printf("Finding relationships for user...\n")
	f, err := followers()
	if err != nil {
		return err
	}
	fmt.Printf("Found %d followers\n", f)

	p := influxdb2.NewPointWithMeasurement("followers").
		AddTag("user", username).
		AddField("value", f).
		SetTime(time.Now())
	return w.WritePoint(context.Background(), p)
}

func init() {
	rootCmd.AddCommand(relationshipsCmd)
}
