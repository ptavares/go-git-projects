package github

import (
	"context"
	"fmt"
	"git-projects/exception"
	"git-projects/git"
	"git-projects/helper"

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
	// If domain is not set -> use default client
	if domain == "" {
		client = github.NewClient(tc)
	} else {
		// Else, it's a EnterpriseClient
		client, err = github.NewEnterpriseClient(domain, domain, tc)
		if err != nil {
			return &Client{}, exception.InitGitHubClientError(err)
		}
	}
	// Test connexion
	user, resp, err := client.Users.Get(ctx, "")
	if err != nil {
		return &Client{}, exception.InitGitHubClientError(err)
	}

	// Rate.Limit should most likely be 5000 when authorized.
	logger.Infof("Rate Limit : %#v\n", resp.Rate)
	// If a Token Expiration has been set, it will be displayed.
	if !resp.TokenExpiration.IsZero() {
		logger.Infof("Token Expiration: %v\n", resp.TokenExpiration)
	}

	logger.Debugf("\n%v\n", github.Stringify(user))

	return &Client{
		client: client,
		listOptions: &github.ListOptions{
			PerPage: 20,
			Page:    1,
		}}, nil
}

// GetProjectsFromGID : Retrieve all user Github's Group projects
func (c *Client) GetProjectsFromGID(gid string, isUser bool) ([]*git.Group, error) {
	logger := helper.GetLogger()

	logger.With(
		zap.String("gid", gid),
		zap.Bool("isUser", isUser),
	).Debug("start - GetProjectsFromGID")

	groupMap := make(map[string]*git.Group)
	var resp *github.Response
	var repositories []*github.Repository
	var err error
	var options interface{}

	if isUser {
		options = &github.RepositoryListOptions{
			Type:        "all",
			Sort:        "created",
			Direction:   "asc",
			ListOptions: *c.listOptions,
		}
	} else {
		options = &github.RepositoryListByOrgOptions{
			Type:        "all",
			Sort:        "created",
			Direction:   "asc",
			ListOptions: *c.listOptions,
		}
	}

	for {
		if isUser {
			repositories, resp, err = c.client.Repositories.List(context.Background(), gid, options.(*github.RepositoryListOptions))
		} else {
			repositories, resp, err = c.client.Repositories.ListByOrg(context.Background(), gid, options.(*github.RepositoryListByOrgOptions))
		}

		if err != nil {
			return []*git.Group{},
				exception.SimpleError(
					fmt.Sprintf("failed to retrieve Github Projects. Is '%s' a %s ? ", gid, helper.Ternary(isUser, "user", "organization")), err)
		}

		for _, repository := range repositories {
			zap.S().With(
				zap.String("path", *repository.Owner.Login),
				zap.String("Name", repository.GetName()),
			).Debug("Repository Info")
			if group, ok := groupMap[*repository.Owner.Login]; ok {
				group.Projects = append(group.Projects, &git.Project{
					Name:          repository.GetName(),
					SSHURLToRepo:  repository.GetSSHURL(),
					HTTPURLToRepo: repository.GetHTMLURL(),
				})
			} else {
				projectList := make([]*git.Project, 0)
				projectList = append(projectList, &git.Project{
					Name:          repository.GetName(),
					SSHURLToRepo:  repository.GetSSHURL(),
					HTTPURLToRepo: repository.GetHTMLURL(),
				})
				groupMap[*repository.Owner.Login] = &git.Group{
					Path:     *repository.Owner.Login,
					Projects: projectList}
			}
		}

		// Exit the loop when we've seen all pages.
		if resp.NextPage == 0 {
			break
		}

		// Update the page number to get the next page.
		if isUser {
			options.(*github.RepositoryListOptions).Page = resp.NextPage
		} else {
			options.(*github.RepositoryListByOrgOptions).Page = resp.NextPage
		}
	}
	values := []*git.Group{}
	for _, value := range groupMap {
		values = append(values, value)
	}

	return values, nil
}
