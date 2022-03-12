package git

import (
	"fmt"
	"git-projects/exception"
	"git-projects/helper"
	"os"
	"path/filepath"

	gitv5 "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

// NewTokenGitAuthentication : Create a Git Token authentication
func NewTokenGitAuthentication(token string) *Auth {
	auth := &Auth{}
	auth.AuthType = Basic
	auth.BasicAuth = &http.BasicAuth{
		Username: "",
		Password: token,
	}
	return auth
}

// NewBasicGitAuthentication : Create a Git Basic authentication
func NewBasicGitAuthentication(userName, password string) *Auth {
	auth := &Auth{}
	auth.AuthType = Basic
	auth.BasicAuth = &http.BasicAuth{
		Username: userName,
		Password: password,
	}
	return auth
}

// NewSSHGitAuthentication : create a Git SSH authentication
func NewSSHGitAuthentication(privateKeyFilePath, password string) (*Auth, error) {
	logger := helper.GetLogger()
	auth := &Auth{}
	logger.Debug("Check if private key file exist")
	_, err := os.Stat(privateKeyFilePath)
	if err != nil {
		return auth, exception.ReadPrivateKeyFileError(privateKeyFilePath, err)
	}
	logger.Debug("Try to generrate public keys")
	publicKeys, err := ssh.NewPublicKeysFromFile("git", privateKeyFilePath, password)
	if err != nil {
		return auth, exception.GeneratePublicKeysError(err)
	}
	auth.AuthType = SSH
	auth.SSHAuth = publicKeys
	return auth, nil
}

// CloneGroupProjects: Clone projets
func (a *Auth) CloneGroupProjects(groups []*Group, destination string) error {
	logger := helper.GetLogger()
	var err error

	// Set destination as current directory if not defined
	if destination == "" {
		logger.Debug("No destination directory specified, will use current directory")
		destination, err = os.Getwd()
		if err != nil {
			return exception.GetPWDDirError(err)
		}
	}
	// Iterate over all groups
	for _, group := range groups {
		logger.Infof("--> Found %d repositories in %s", len(group.Projects), group.Path)
		// Create DirTree matching Group Path
		subDir := filepath.FromSlash(fmt.Sprintf("%s/%s", destination, group.Path))
		if err := os.MkdirAll(subDir, 0755); err != nil {
			return exception.CreateDirError(subDir, err)
		}
		// Clone each projects if not already exist
		for _, project := range group.Projects {
			// Create clone directory directory using group path and project name
			cloneDir := filepath.FromSlash(fmt.Sprintf("%s/%s", subDir, project.Name))
			// Already exist dir -> not clone to avoid error
			if _, err := os.Stat(cloneDir); !os.IsNotExist(err) {
				logger.Infof("-----> Skipping %s, directory %s already exist", project.Name, cloneDir)
				continue
			}
			// Else clone project using wanted auth type (SSH or HTTP)
			logger.Infof("-----> Clonning %s into directory %s\n", project.Name, cloneDir)
			commitInfo := a.cloneProject(cloneDir, project)
			// Print result of clone with last commit information
			if commitInfo != "" {
				logger.Infof("%s", commitInfo)
			}
		}
	}
	return nil
}

// cloneProject : execute a clone of project in cloneDir
func (a *Auth) cloneProject(cloneDir string, project *Project) string {
	logger := helper.GetLogger()

	// Define authentication and URL
	var auth interface{}
	var url string
	switch a.AuthType {
	case Basic:
		auth = a.BasicAuth
		url = project.HTTPURLToRepo
	case SSH:
		auth = a.SSHAuth
		url = project.SSHURLToRepo
	}

	r, err := gitv5.PlainClone(cloneDir, false, &gitv5.CloneOptions{
		Auth:     auth.(transport.AuthMethod),
		URL:      url,
		Progress: os.Stdout,
	})

	if err != nil {
		logger.Errorf("cloning repository %s error : %s", project.HTTPURLToRepo, err)
		return ""
	}
	// Retrieving the branch being pointed by HEAD
	ref, err := r.Head()
	if err != nil {
		logger.Errorf("getting repository HEAD branch for %s error : %s", project.HTTPURLToRepo, err)
		return ""
	}
	// Retrieving the commit object
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		logger.Errorf("getting last commit object for %s error : %s", project.HTTPURLToRepo, err)
		return ""
	}
	return commit.String()
}
