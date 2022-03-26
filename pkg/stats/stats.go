package stats

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/go-github/v42/github"
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
	participantStats map[int]UserStats
}

// CheckArgs parses the cli input before we run our program
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

// Run is the entry point into the stats package, this is where we will begin to
// make our api requests to collect pull request data, and then parse the results
// and return our user stats information collected.
func Run(ctx context.Context, options *UserStatsOptions) error {
	manager, err := newManager(options)
	if err != nil {
		return err
	}

	prList, err := manager.pullRequestList(ctx)
	if err != nil {
		return err
	}

	err = manager.parsePullRequestList(ctx, prList)
	if err != nil {
		return err
	}

	return nil
}

// PullRequestList looks for ALL pull requests that occur AFTER
// the managers options AfterDate and returns a list of PR numbers
func (m *StatsManager) pullRequestList(ctx context.Context) ([]int, error) {
	l, _, err := m.ghcli.PullRequests.List(ctx, m.options.Owner, m.options.Repo, &m.options.ListOptions)
	if err != nil {
		return nil, err
	}

	var n []int
	for _, v := range l {
		if v.CreatedAt.After(m.options.AfterDate) {
			continue
		}
		n = append(n, *v.Number)
	}
	return n, nil
}

func newManager(o *UserStatsOptions) (*StatsManager, error) {
	var m = new(StatsManager)

	m.ghcli = github.NewClient(nil)
	m.options = o
	m.options.ListOptions = github.PullRequestListOptions{
		State: "all",
	}

	return m, nil
}
