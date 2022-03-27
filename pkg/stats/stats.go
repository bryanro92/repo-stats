package stats

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
)

// Run is the entry point into the stats package, this is where we will begin to
// make our api requests to collect pull request data, and then parse the results
// and return our user stats information collected.
func Run(ctx context.Context, options *UserStatsOptions) error {
	manager, err := newManager(ctx, options)
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
	manager.printResults()
	return nil
}

func (m *StatsManager) printResults() {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "User\tApproved\tComments\tChanges Requested\tTotal Interactions\tTotal PRs\tPR List")
	for _, u := range m.participantStats {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\n", u.Username, u.Approvals, u.Comments, u.ChangesRequested, u.total(), u.totalPRs(), u.uniquePRs())
		// fmt.Printf("user: %v approve: %v comment: %v: changes: %v total: %v prNumReviewed: %v prNumList: %v\n", u.Username, u.Approvals, u.Comments, u.ChangesRequested, u.total(), len(u.PullList), u.PullList)
	}
	w.Flush()
}

func (m *UserStats) total() int {
	return m.Approvals + m.ChangesRequested + m.Comments
}

func (m *UserStats) totalPRs() int {
	return len(m.uniquePRs())
}

func (m *UserStats) uniquePRs() []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range m.PullList {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
