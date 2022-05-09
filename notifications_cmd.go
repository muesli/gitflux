package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v32/github"
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

	total := 0
	unread := 0

	opt := &github.NotificationListOptions{
		ListOptions: github.ListOptions{
			PerPage: 50,
		},
	}
	for {
		notifications, resp, err := clientv3.Activity.ListNotifications(context.Background(), opt)
		if err != nil {
			return err
		}

		for _, v := range notifications {
			total++
			if *v.Unread {
				unread++
			}
		}

		if resp.NextPage == 0 || len(notifications) == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	fmt.Printf("Found %d total, %d unread notifications\n", total, unread)

	p := influxdb2.NewPointWithMeasurement("notifications").
		AddTag("user", username).
		AddField("value", unread).
		SetTime(time.Now())
	return w.WritePoint(context.Background(), p)
}

func init() {
	rootCmd.AddCommand(notificationsCmd)
}
