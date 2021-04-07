package main

import (
	"context"
	"strings"
	"time"

	"github.com/shurcooL/githubv4"
)

type QLUser struct {
	Login     githubv4.String
	Name      githubv4.String
	AvatarURL githubv4.String
	URL       githubv4.String
}

type User struct {
	Login     string
	Name      string
	AvatarURL string
	URL       string
}

var viewerQuery struct {
	Viewer struct {
		Login githubv4.String
	}
}

func getUsername() (string, error) {
	err := client.Query(context.Background(), &viewerQuery, nil)
	if err != nil {
		if strings.Contains(err.Error(), "abuse-rate-limits") {
			time.Sleep(time.Minute)
			return getUsername()
		}
		return "", err
	}

	return string(viewerQuery.Viewer.Login), nil
}

var followersQuery struct {
	User struct {
		Login     githubv4.String
		Followers struct {
			TotalCount githubv4.Int
		}
	} `graphql:"user(login:$username)"`
}

func followers() (int, error) {
	variables := map[string]interface{}{
		"username": githubv4.String(username),
	}
	err := client.Query(context.Background(), &followersQuery, variables)
	if err != nil {
		if strings.Contains(err.Error(), "abuse-rate-limits") {
			time.Sleep(time.Minute)
			return followers()
		}
		return 0, err
	}

	return int(followersQuery.User.Followers.TotalCount), nil
}

func UserFromQL(user QLUser) User {
	return User{
		Login:     string(user.Login),
		Name:      string(user.Name),
		AvatarURL: string(user.AvatarURL),
		URL:       string(user.URL),
	}
}
