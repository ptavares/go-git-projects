package cmd

import (
	"fmt"
	"git-projects/config"
	"git-projects/git"
	"git-projects/gitlab"
	"git-projects/helper"
	"git-projects/internal/name"
	"os"

	"github.com/spf13/cobra"
)

// gitlabCloneCmd represents the clone command
var gitlabCloneCmd = &cobra.Command{
	Use:   "clone",
	Short: fmt.Sprintf("Perform %s Gitlab clone actions", name.ApplicationName),
	Long: fmt.Sprintf(`
=======================================================================
=                     %s gitlab clone                       =
=======================================================================

Command to clone Gitlab projects

`, name.ApplicationName),
}

// gitlabCloneHTTPCmd represents the clone HTTP command
var gitlabCloneHTTPCmd = &cobra.Command{
	Use:   "http",
	Short: fmt.Sprintf("Perform %s Gitlab HTTP clone actions", name.ApplicationName),
	Long: fmt.Sprintf(`
=======================================================================
=                   %s gitlab clone HTTP                    =
=======================================================================

Command to clone Gitlab projects from repositories's HTTP URL

`, name.ApplicationName),
	PreRun: checkGitlabCloneArgument,
	Run:    executeGitlabCloneHTTP,
}

// gitlabCloneSSHCmd represents the clone SSH command
var gitlabCloneSSHCmd = &cobra.Command{
	Use:   "ssh",
	Short: fmt.Sprintf("Perform %s Gitlab SSH clone actions", name.ApplicationName),
	Long: fmt.Sprintf(`
=======================================================================
=                   %s gitlab clone SSH                   =
=======================================================================

Command to clone Gitlab projects from repositories's SSH URL

`, name.ApplicationName),
	PreRun: checkGitlabCloneArgument,
	Run:    executeGitlabCloneSSH,
}

// init : Init Gitlab Clone sub commands
func init() {
	gitlabCmd.AddCommand(gitlabCloneCmd)
	gitlabCloneCmd.AddCommand(gitlabCloneHTTPCmd)
	gitlabCloneCmd.AddCommand(gitlabCloneSSHCmd)

	// -> Flags for Basic Auth
	gitlabCloneHTTPCmd.PersistentFlags().StringVarP(&basicAuthUsername, "basic-auth-username", "", "", fmt.Sprintf("username to use to clone repository throw HTTP URL <%s_%s>", config.ENV_PREFIX, "BASIC_AUTH_USR"))
	gitlabCloneHTTPCmd.PersistentFlags().StringVarP(&basicAuthPassword, "basic-auth-password", "", "", fmt.Sprintf("password related to 'basic-auth-username' <%s_%s>", config.ENV_PREFIX, "BASIC_AUTH_PWD"))
	gitlabCloneHTTPCmd.PersistentFlags().StringVarP(&basicAuthToken, "basic-auth-token", "", "", fmt.Sprintf("token to use to clone repository throw HTTP URL if different from 'api-token' <%s_%s>", config.ENV_PREFIX, "BASIC_AUTH_TOKEN"))

	// -> Flags for SSH
	gitlabCloneSSHCmd.PersistentFlags().StringVarP(&sshPrivateKeyPath, "ssh-private-key-path", "", "", fmt.Sprintf("path to private key file used to clone repository throw SSH URL <%s_%s>", config.ENV_PREFIX, "SSH_KEY_PATH"))
	gitlabCloneSSHCmd.PersistentFlags().StringVarP(&sshPrivateKeyPwd, "ssh-private-key-password", "", "", fmt.Sprintf("optional password to decrypt private key <%s_%s>", config.ENV_PREFIX, "SSH_KEY_PWD"))

	// Add EnvName Param to config
	config.AddEnvParam("BASIC_AUTH_USR")
	config.AddEnvParam("BASIC_AUTH_PWD")
	config.AddEnvParam("BASIC_AUTH_TOKEN")
	config.AddEnvParam("SSH_KEY_PATH")
	config.AddEnvParam("SSH_KEY_PWD")

}

func checkGitlabCloneArgument(cmd *cobra.Command, args []string) {
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

// getAllProjects : Retrieve all Gitlab projects (optional : for a given group)
func getAllProjects() []*git.Group {
	// Create new client
	client, err := gitlab.NewClient(apiUserToken, fmt.Sprintf("https://%s", baseDomain))
	if err != nil {
		helper.HandleErrorExit(err)
	}

	// Retrieve all projects
	// 491 for cld
	groups, err := client.GetProjectsFromGID(gid)
	if err != nil {
		helper.HandleErrorExit(err)
	}
	logger.Debugf("Groups : %v", groups)

	return groups
}

// executeGitlabCloneHTTP : clone gitlab projects using HTTP protocol
func executeGitlabCloneHTTP(cmd *cobra.Command, args []string) {

	// Create Gitlab authentication
	var auth *git.Auth
	if basicAuthUsername != "" {
		auth = git.NewBasicGitAuthentication(basicAuthUsername, basicAuthPassword)
	} else {
		auth = git.NewTokenGitAuthentication(helper.Ternary(basicAuthToken != "", basicAuthToken, apiUserToken))
	}

	// Get all projects by Group
	projectsByGroup := getAllProjects()
	// Clone projects
	err := auth.CloneGroupProjects(projectsByGroup, destination)
	if err != nil {
		helper.HandleErrorExit(err)
	}
}

// executeGitlabCloneSSH : clone gitlab projects using SSH protocol
func executeGitlabCloneSSH(cmd *cobra.Command, args []string) {
	// Create Gitlab authentication
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
