package stats

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

// PullRequestList looks for ALL pull requests that occur AFTER
// the managers options AfterDate and returns a list of PR numbers
func (m *StatsManager) pullRequestList(ctx context.Context) ([]int, error) {
	l, _, err := m.ghcli.PullRequests.List(ctx, m.options.Owner, m.options.Repo, &m.options.ListOptions)
	if err != nil {
		return nil, err
	}

	var n []int
	for _, v := range l {
		if *v.Number < 2020 {
			// if v.CreatedAt.After(m.options.AfterDate) {
			continue
		}
		n = append(n, *v.Number)
	}
	return n, nil
}

// ParsePullRequestList loops through our list of pr numbers,
// and requests the manager parse pr reviews
func (m *StatsManager) parsePullRequestList(ctx context.Context, prList []int) error {
	for _, n := range prList {
		err := m.parsePullRequestReviews(ctx, m.options.Owner, m.options.Repo, n)
		if err != nil {
			fmt.Println("err parsing pr:", err)
		}
	}
	return nil
}

// ParsePullRequestReviews analyzes the activity of pr reviews
// and updates our manager struct with the data
func (m *StatsManager) parsePullRequestReviews(ctx context.Context, owner string, repo string, number int) error {
	cm, _, err := m.ghcli.PullRequests.ListReviews(ctx, m.options.Owner, m.options.Repo, number, &github.ListOptions{})
	if err != nil {
		return err
	}
	if len(cm) == 0 {
		fmt.Println("No reviews yet")
	}
	for _, v := range cm {
		_, ok := m.participantStats[*v.User.ID]
		if !ok {
			m.participantStats[*v.User.ID] = &UserStats{
				Username: *v.User.Login,
			}
		}
		switch *v.State {
		case "APPROVED":
			m.participantStats[*v.User.ID].Approvals += 1
		case "COMMENTED":
			m.participantStats[*v.User.ID].Comments += 1
		case "CHANGES_REQUESTED":
			m.participantStats[*v.User.ID].ChangesRequested += 1
		default:
			fmt.Println(*v.State)
		}
	}
	fmt.Println("------------------------------------------")
	fmt.Println("PR number:", number)
	for user, num := range m.participantStats {
		fmt.Println("User id:", user, "name:", num.Username, " - contributions: ", num)
	}
	return err
}
