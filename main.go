package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v43/github"
)

type PRStats struct {
	PullsOpened      int
	PullsMerged      int
	Approvals        int
	Comments         int
	ChangesRequested int
}

func main() {
	ctx := context.Background()
	client := github.NewClient(nil)
	var u = make(map[string]int)
	//options := &github.PullRequestListOptions{
	//	State: "all",
	//}
	cm, _, err := client.PullRequests.ListReviews(ctx, "Azure", "ARO-RP", 2009, &github.ListOptions{})
	//cm, resp, err := client.Repositories.ListComments(ctx, "Azure", "ARO-RP", &github.ListOptions{})
	//pr, resp, err := client.PullRequests.List(ctx, "Azure", "ARO-RP", options)
	if err != nil {
		fmt.Errorf("err %v", err)
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
}
