package github

import (
	"context"
	"fmt"
	"git-projects/helper"
	"log"

	"github.com/google/go-github/v43/github"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

// NewClient: Create a Gitlab client
func NewClient(token, domain string) (*Client, error) {
	logger := helper.GetLogger()

	logger.With(
		zap.String("token", helper.MaskPassword(token)),
		zap.String("domain", domain),
	).Debug("start - NewClient")

	// Setup new Github client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	var client *github.Client
	var err error
	if domain == "" {
		client = github.NewClient(tc)
	} else {
		client, err = github.NewEnterpriseClient(domain, domain, tc)
	}
	user, resp, err := client.Users.Get(ctx, "")
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return &Client{}, err
	}

	// Rate.Limit should most likely be 5000 when authorized.
	log.Printf("Rate: %#v\n", resp.Rate)
	// If a Token Expiration has been set, it will be displayed.
	if !resp.TokenExpiration.IsZero() {
		log.Printf("Token Expiration: %v\n", resp.TokenExpiration)
	}

	fmt.Printf("\n%v\n", github.Stringify(user))

	return &Client{client: client}, nil
}
