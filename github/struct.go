package github

import "github.com/google/go-github/v43/github"

// A Client manage communication with the Github API
type Client struct {
	client      *github.Client
	listOptions *github.ListOptions
}
