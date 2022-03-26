package stats

import (
	"context"
	"fmt"

	"github.com/google/go-github/v42/github"
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
		if v.CreatedAt.After(m.options.AfterDate) {
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
	return fmt.Errorf("error occuured during parsing")
}

// ParsePullRequestReviews analyzes the activity of pr reviews
// and updates our manager struct with the data
func (m *StatsManager) parsePullRequestReviews(ctx context.Context, owner string, repo string, number int) error {
	// var contributorStats = make(map[int]UserStats)

	var u = make(map[string]int)
	//options := &github.PullRequestListOptions{
	//	State: "all",
	//}
	cm, _, err := m.ghcli.PullRequests.ListReviews(ctx, m.options.Owner, m.options.Repo, number, &github.ListOptions{})
	//cm, resp, err := client.Repositories.ListComments(ctx, "Azure", "ARO-RP", &github.ListOptions{})
	//pr, resp, err := client.PullRequests.List(ctx, "Azure", "ARO-RP", options)
	if err != nil {
		return err
	}
	if len(cm) == 0 {
		fmt.Println("No reviews yet")
	}
	c := 0
	for k, v := range cm {
		fmt.Println(k)
		fmt.Println(*v.State)
		fmt.Println(*v.User.Login)
		u[*v.User.Login]++
		c++
	}
	fmt.Println("------------------------------------------")
	fmt.Println(cm)
	fmt.Println("------------------------------------------")
	fmt.Println("PR number: 2009")
	for user, num := range u {
		fmt.Println("User:", user, " - contributions: ", num)
	}
	return err
}
