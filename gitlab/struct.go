package gitlab

import (
	"github.com/xanzy/go-gitlab"
)

// A Client manage communication with the Gitlab API
type Client struct {
	client       *gitlab.Client
	sort         string
	orderBy      string
	topLevelOnly bool
	listOptions  gitlab.ListOptions
}
