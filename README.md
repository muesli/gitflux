# gitflux

Track your GitHub projects in InfluxDB and create beautiful graphs with Grafana

## Features

Lets you track these things:

- Yourself
  - [x] Follower counts
  - [x] Notifications
- Repositories
  - [x] Stars
  - [x] Forks
  - [x] Watchers
  - [x] Commits
- Issues
  - [x] State
  - [x] Assignees
  - [x] Labels
- PRs
  - [x] State
  - [x] Assignees
  - [x] Labels

## Usage

Import statistics for all your source repositories:

```
$ gitflux repository
Finding user's source repos...
Found 83 repos
Parsing muesli/gitflux
    Finding PRs for repo...
    Found 38 PRs!
    Finding issues for repo...
    Found 39 issues!
Parsing muesli/duf
...
```

Import statistics for a specific repository:

```
$ gitflux repository muesli/gitflux
Parsing muesli/gitflux
    Finding PRs for repo...
    Found 38 PRs!
    Finding issues for repo...
    Found 39 issues!
```

Import relationship statistics:

```
$ gitflux relationships
Finding relationships for user...
Found 1109 followers
```

Import notification statistics:

```
$ gitflux notifications
Finding notifications for user...
Found 14 unread notifications
```

### Flags

```
--influx string          InfluxDB address (default "http://localhost:8086")
--influx-bucket string   InfluxDB bucket (default "github")
--influx-token string    InfluxDB auth token
```

## Screenshots

### Graphs about you

![followers](/screenshots/user_followers.png)
![notifications](/screenshots/user_notifications.png)

### Graphs about all your source repos

![stars](/screenshots/repo_stars.png)
![forks](/screenshots/repo_forks.png)
![watchers](/screenshots/repo_watchers.png)
![commits](/screenshots/repo_commits.png)
![issues](/screenshots/repo_issues.png)
![prs](/screenshots/repo_prs.png)

### Graphs about individual projects

![stars](/screenshots/project_stars.png)
![forks](/screenshots/project_forks.png)
![watchers](/screenshots/project_watchers.png)
![commits](/screenshots/project_commits.png)
![issues](/screenshots/project_issues.png)
![issue labels](/screenshots/project_issues_labels.png)
![issue bars](/screenshots/project_issues_labels_bars.png)
![prs](/screenshots/project_prs.png)
![pr labels](/screenshots/project_prs_labels.png)
![pr bars](/screenshots/project_prs_labels_bars.png)

## TODOs

- Add a `docker-compose.yml` with the following services:
  - InfluxDB
  - Grafana
  - gitflux
- More graphs?
