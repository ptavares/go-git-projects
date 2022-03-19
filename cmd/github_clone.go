package cmd

import (
	"fmt"
	"git-projects/config"
	"git-projects/git"
	"git-projects/github"
	"git-projects/helper"
	"git-projects/internal/name"
	"os"

	"github.com/spf13/cobra"
)

// githubCloneCmd represents the clone command
var githubCloneCmd = &cobra.Command{
	Use:   "clone",
	Short: fmt.Sprintf("Perform %s Github clone actions", name.ApplicationName),
	Long: fmt.Sprintf(`
=======================================================================
=                     %s github clone                       =
=======================================================================

Command to clone Github projects

`, name.ApplicationName),
}

// githubCloneHTTPCmd represents the clone HTTP command
var githubCloneHTTPCmd = &cobra.Command{
	Use:   "http",
	Short: fmt.Sprintf("Perform %s Github HTTP clone actions", name.ApplicationName),
	Long: fmt.Sprintf(`
=======================================================================
=                   %s github clone HTTP                    =
=======================================================================

Command to clone Github projects from repositories's HTTP URL

`, name.ApplicationName),
	PreRun: checkGithubCloneArgument,
	Run:    executeGithubCloneHTTP,
}

// githubCloneSSHCmd represents the clone SSH command
var githubCloneSSHCmd = &cobra.Command{
	Use:   "ssh",
	Short: fmt.Sprintf("Perform %s Github SSH clone actions", name.ApplicationName),
	Long: fmt.Sprintf(`
=======================================================================
=                   %s github clone SSH                   =
=======================================================================

Command to clone Github projects from repositories's SSH URL

`, name.ApplicationName),
	PreRun: checkGithubCloneArgument,
	Run:    executeGithubCloneSSH,
}

// init : Init Github Clone sub commands
func init() {
	githubCmd.AddCommand(githubCloneCmd)
	githubCloneCmd.AddCommand(githubCloneHTTPCmd)
	githubCloneCmd.AddCommand(githubCloneSSHCmd)

	// -> Flags for Basic Auth
	githubCloneHTTPCmd.PersistentFlags().StringVarP(&basicAuthUsername, "basic-auth-username", "", "", fmt.Sprintf("username to use to clone repository throw HTTP URL <%s_%s>", config.ENV_PREFIX, "BASIC_AUTH_USR"))
	githubCloneHTTPCmd.PersistentFlags().StringVarP(&basicAuthPassword, "basic-auth-password", "", "", fmt.Sprintf("password related to 'basic-auth-username' <%s_%s>", config.ENV_PREFIX, "BASIC_AUTH_PWD"))
	githubCloneHTTPCmd.PersistentFlags().StringVarP(&basicAuthToken, "basic-auth-token", "", "", fmt.Sprintf("token to use to clone repository throw HTTP URL if different from 'api-token' <%s_%s>", config.ENV_PREFIX, "BASIC_AUTH_TOKEN"))

	// -> Flags for SSH
	githubCloneSSHCmd.PersistentFlags().StringVarP(&sshPrivateKeyPath, "ssh-private-key-path", "", "", fmt.Sprintf("path to private key file used to clone repository throw SSH URL <%s_%s>", config.ENV_PREFIX, "SSH_KEY_PATH"))
	githubCloneSSHCmd.PersistentFlags().StringVarP(&sshPrivateKeyPwd, "ssh-private-key-password", "", "", fmt.Sprintf("optional password to decrypt private key <%s_%s>", config.ENV_PREFIX, "SSH_KEY_PWD"))

	// Add EnvName Param to config
	config.AddEnvParam("BASIC_AUTH_USR")
	config.AddEnvParam("BASIC_AUTH_PWD")
	config.AddEnvParam("BASIC_AUTH_TOKEN")
	config.AddEnvParam("SSH_KEY_PATH")
	config.AddEnvParam("SSH_KEY_PWD")

}

func checkGithubCloneArgument(cmd *cobra.Command, args []string) {
	// Basic auth
	fillStringParam("basic_auth_usr", config.GetConfig().BasicAuth.UserName, &basicAuthUsername)
	fillStringParam("basic_auth_pwd", config.GetConfig().BasicAuth.Password, &basicAuthPassword)
	fillStringParam("basic_auth_token", config.GetConfig().BasicAuth.Token, &basicAuthToken)
	// SSH Auth
	fillStringParam("ssh_key_path", config.GetConfig().SSHAuth.PrivateKeyPath, &sshPrivateKeyPath)
	fillStringParam("ssh_key_pwd", config.GetConfig().SSHAuth.PrivateKeyPassword, &sshPrivateKeyPwd)

	// Check authentication type
	if sshPrivateKeyPath != "" && (basicAuthUsername != "" || basicAuthToken != "") {
		logger.Info("error : only one type of authentication is available, choose between ssh or basic authentication ")
		os.Exit(1)
	}
	if basicAuthUsername != "" && basicAuthToken != "" {
		logger.Info("error : only one type of basic authentication is available, choose between username/password or token")
		os.Exit(1)
	}
}

// getAllProjects : Retrieve all Github projects (optional : for a given group)
func getAllGithubProjects() []*git.Group {
	// Create new client
	domain := helper.Ternary(baseDomain == default_github_domain, "", baseDomain)
	client, err := github.NewClient(apiUserToken, domain)
	if err != nil {
		helper.HandleErrorExit(err)
	}

	// Retrieve all projects
	// 491 for cld
	groups, err := client.GetProjectsFromGID(gid, isUser)
	if err != nil {
		helper.HandleErrorExit(err)
	}
	logger.Debugf("Groups : %v", groups)

	return groups
}

// executeGithubCloneHTTP : clone github projects using HTTP protocol
func executeGithubCloneHTTP(cmd *cobra.Command, args []string) {

	// Create Github authentication
	var auth *git.Auth
	if basicAuthUsername != "" {
		auth = git.NewBasicGitAuthentication(basicAuthUsername, basicAuthPassword)
	} else {
		auth = git.NewTokenGitAuthentication(helper.Ternary(basicAuthToken != "", basicAuthToken, apiUserToken))
	}

	// Get all projects by Group
	projectsByGroup := getAllGithubProjects()
	// Clone projects
	err := auth.CloneGroupProjects(projectsByGroup, destination)
	if err != nil {
		helper.HandleErrorExit(err)
	}
}

// executeGithubCloneSSH : clone github projects using SSH protocol
func executeGithubCloneSSH(cmd *cobra.Command, args []string) {
	// Create Github authentication
	auth, err := git.NewSSHGitAuthentication(sshPrivateKeyPath, sshPrivateKeyPwd)
	if err != nil {
		helper.HandleErrorExit(err)
	}

	// Get all projects by Group
	projectsByGroup := getAllProjects()

	// Clone projects
	err = auth.CloneGroupProjects(projectsByGroup, destination)
	if err != nil {
		helper.HandleErrorExit(err)
	}
}
