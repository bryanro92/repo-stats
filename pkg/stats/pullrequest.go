package stats

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

// PullRequestList looks for ALL pull requests that occur AFTER
// the managers options AfterDate and returns a list of PR numbers
func (m *StatsManager) pullRequestList(ctx context.Context) ([]int, error) {
	// TODO handle pagination
	// It looks like we get 30 PR's per page
	var n []int

	// We will get 8 pages of results, 8*30= ~240 PRs
	// Each PR requires an additional API call to collect the reviews
	// This gives us about 248 API Calls
	for m.options.ListOptions.Page < 9 {
		l, resp, err := m.ghcli.PullRequests.List(ctx, m.options.Owner, m.options.Repo, &m.options.ListOptions)
		if resp.NextPage != 0 {
			m.options.ListOptions.ListOptions.Page = resp.NextPage
		}

		if err != nil {
			return nil, err
		}
		for _, v := range l {
			if *v.Number < 1400 {
				// if v.CreatedAt.After(m.options.AfterDate) {
				continue
			}
			_, ok := m.participantStats[*v.User.ID]
			if !ok {
				m.participantStats[*v.User.ID] = &UserStats{
					Username: *v.User.Login,
				}
			}
			n = append(n, *v.Number)
		}
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
		case "PENDING":
		case "DISMISSED":
		default:
			fmt.Println(*v.State)
		}
		m.participantStats[*v.User.ID].PullList = append(m.participantStats[*v.User.ID].PullList, number)
	}
	return nil
}
