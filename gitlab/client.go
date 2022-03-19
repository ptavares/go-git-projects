package gitlab

import (
	"fmt"
	"git-projects/exception"
	"git-projects/git"
	"git-projects/helper"
	"net/url"
	"sort"

	"github.com/xanzy/go-gitlab"
	"go.uber.org/zap"
)

// NewClient: Create a Gitlab client
func NewClient(token, domain string) (*Client, error) {
	logger := helper.GetLogger()

	logger.With(
		zap.String("token", helper.MaskPassword(token)),
		zap.String("domain", domain),
	).Debug("start - NewClient")

	// Setup new Gitlab client
	git, err := gitlab.NewClient(token, gitlab.WithBaseURL(domain))
	if err != nil {
		return &Client{}, exception.InitGitlabClientError(err)
	}

	// Try to call a simple API request to check if token is OK
	version, _, err := git.Version.GetVersion()
	logger.Debug(version)
	if err != nil {
		switch v := err.(type) {
		case *gitlab.ErrorResponse:
			return &Client{},
				exception.APICallError(fmt.Sprintf("unable to call Gitlab API, status code %d", v.Response.StatusCode))
		case *url.Error:
			return &Client{},
				exception.FetchingURLError(v.URL, v.Err)
		default:
			return &Client{},
				exception.SimpleError("unable to call Gitlab API", err)
		}
	}

	return &Client{
		client: git,
		listOptions: gitlab.ListOptions{
			PerPage: 20,
			Page:    1,
		},
		sort:         "desc",
		orderBy:      "path",
		topLevelOnly: false,
	}, nil
}

// GetProjectsFromGID : Retrieve all user Gitlab's Group with projects
func (c *Client) GetProjectsFromGID(gid string) ([]*git.Group, error) {
	logger := helper.GetLogger()

	logger.With(
		zap.String("gid", gid),
	).Debug("start - GetProjectsFromGID")

	groupList := make([]*git.Group, 0)
	var err error
	var resp *gitlab.Response
	var groups []*gitlab.Group
	var options interface{}

	if gid == "" {
		logger.Debug("List Group without GID")
		options = &gitlab.ListGroupsOptions{
			Sort:         gitlab.String(c.sort),
			OrderBy:      gitlab.String(c.orderBy),
			TopLevelOnly: gitlab.Bool(c.topLevelOnly),
			ListOptions:  c.listOptions,
		}
	} else {
		logger.With(zap.String("gid", gid)).Debug("List Group with GID")
		options = &gitlab.ListDescendantGroupsOptions{
			Sort:        gitlab.String(c.sort),
			OrderBy:     gitlab.String(c.orderBy),
			ListOptions: c.listOptions,
		}
	}

	for {
		if gid == "" {
			groups, resp, err = c.client.Groups.ListGroups(options.(*gitlab.ListGroupsOptions))
		} else {
			groups, resp, err = c.client.Groups.ListDescendantGroups(gid, options.(*gitlab.ListDescendantGroupsOptions))
		}

		if err != nil {
			return []*git.Group{},
				exception.SimpleError("failed to retrieve Gitlab Groups", err)
		}

		for _, group := range groups {
			zap.S().With(
				zap.String("path", group.FullPath),
				zap.Int("gid", group.ID),
			).Debug("Group Info")
			groupList = append(groupList, &git.Group{Path: group.FullPath, GID: group.ID})
		}

		// Exit the loop when we've seen all pages.
		if resp.NextPage == 0 {
			break
		}

		// Update the page number to get the next page.
		if gid == "" {
			options.(*gitlab.ListGroupsOptions).Page = resp.NextPage
		} else {
			options.(*gitlab.ListDescendantGroupsOptions).Page = resp.NextPage
		}
	}
	sort.Slice(groupList, func(i, j int) bool {
		return groupList[i].Path < groupList[j].Path
	})

	// Can be a unique group
	if len(groupList) == 0 {
		group, _, err := c.client.Groups.GetGroup(gid, &gitlab.GetGroupOptions{})
		if err != nil {
			return []*git.Group{},
				exception.SimpleError("failed to retrieve Gitlab Groups", err)
		}
		groupList = append(groupList, &git.Group{Path: group.FullPath, GID: group.ID})
	}

	logger.Debug("Get Projects for each group")
	for _, group := range groupList {
		projects, err := c.getProjects(*group)
		if err != nil {
			return groupList, err
		}
		group.Projects = projects
	}

	return groupList, nil
}

// getProjects : Retrieve all group projects
func (c *Client) getProjects(group git.Group) ([]*git.Project, error) {
	logger := helper.GetLogger()
	logger.With(zap.Int("gid", group.GID), zap.String("path", group.Path)).Debug("List Project for Group")

	projectList := make([]*git.Project, 0)

	options := &gitlab.ListGroupProjectsOptions{
		Sort:             gitlab.String(c.sort),
		OrderBy:          gitlab.String(c.orderBy),
		IncludeSubgroups: gitlab.Bool(false),
		ListOptions:      c.listOptions,
	}

	for {
		projects, resp, err := c.client.Groups.ListGroupProjects(group.GID, options)

		if err != nil {
			return []*git.Project{},
				exception.SimpleError(fmt.Sprintf("failed to retrieve projects for group %s", group.Path), err)
		}

		if len(projects) > 0 {
			logger.Debugf("Found %d project(s)", len(projects))
			for _, project := range projects {
				zap.S().With(
					zap.String("Name", project.Name),
					zap.Int("pid", project.ID),
					zap.String("sshUrl", project.SSHURLToRepo),
					zap.String("httpUrl", project.HTTPURLToRepo),
				).Debug("  -> Project Info")
				projectList = append(projectList, &git.Project{
					Name:          project.Name,
					PID:           project.ID,
					SSHURLToRepo:  project.SSHURLToRepo,
					HTTPURLToRepo: project.HTTPURLToRepo,
				})
			}
		} else {
			logger.Debug("No project found")
		}

		// Exit the loop when we've seen all pages.
		if resp.NextPage == 0 {
			break
		}
		options.Page = resp.NextPage
	}

	return projectList, nil
}
