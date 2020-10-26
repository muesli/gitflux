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
	notificationsCmd = &cobra.Command{
		Use:   "notifications",
		Short: "tracks unread notifications",
		RunE: func(cmd *cobra.Command, args []string) error {
			return parseNotifications(influxWriter)
		},
	}
)

func parseNotifications(w api.WriteAPIBlocking) error {
	fmt.Printf("Finding notifications for user...\n")
	notifications, _, err := clientv3.Activity.ListNotifications(context.Background(), nil)
	if err != nil {
		return err
	}

	unread := 0
	for _, v := range notifications {
		if *v.Unread {
			unread++
		}
	}
	fmt.Printf("Found %d total, %d unread notifications\n", len(notifications), unread)

	p := influxdb2.NewPointWithMeasurement("notifications").
		AddTag("user", username).
		AddField("value", unread).
		SetTime(time.Now())
	return w.WritePoint(context.Background(), p)
}

func init() {
	rootCmd.AddCommand(notificationsCmd)
}
