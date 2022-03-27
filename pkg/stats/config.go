// CheckArgs parses the cli input before we run our program
package stats

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

/*
Package stats will hold our data structures and
do the work to populate activity for a given repository
*/

// UserStats represents a users participation
type UserStats struct {
	Approvals        int
	Comments         int
	ChangesRequested int
	PullsOpened      int
	PullsMerged      int
	Username         string
}

// UserStatsOptions represents configuration options for our queries to generate the stats
type UserStatsOptions struct {
	AfterDate   time.Time
	Owner       string
	Repo        string
	ListOptions github.PullRequestListOptions
}

// StatsManager represents our manager struct containing access to required clients
// and data structures holding our participant data
type StatsManager struct {
	ghcli            *github.Client
	options          *UserStatsOptions
	participantStats map[int64]*UserStats
}

func CheckArgs(args []string) (*UserStatsOptions, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("invalid number of args, expected: 1, got: %v", len(args))
	}
	d, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, fmt.Errorf("couldn't parse parameter to int Days: %v", err)
	}
	var options UserStatsOptions
	options.Owner = args[0]
	options.Repo = args[1]
	options.AfterDate = time.Now().AddDate(0, 0, -d)
	return &options, nil
}

// Usage gives examples of how to run the program
func Usage() {
	fmt.Println("Example usage:")
	fmt.Println("./repo-stats owner repo days")
	fmt.Println("./repo-stats azure aro-rp 90")
}

func newManager(ctx context.Context, o *UserStatsOptions) (*StatsManager, error) {
	var m = new(StatsManager)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GH_PAT")},
	)
	tc := oauth2.NewClient(ctx, ts)
	m.participantStats = make(map[int64]*UserStats)
	m.ghcli = github.NewClient(tc)
	m.options = o
	m.options.ListOptions = github.PullRequestListOptions{
		State: "all",
	}

	return m, nil
}
