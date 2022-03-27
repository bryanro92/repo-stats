package stats

import (
	"context"
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

	return nil
}
