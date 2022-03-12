package cmd

import (
	"fmt"
	"git-projects/config"
	"git-projects/internal/name"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// gitlabCmd represents the config command
var gitlabCmd = &cobra.Command{
	Use:   "gitlab",
	Short: fmt.Sprintf("Perform %s Gitlab actions", name.ApplicationName),
	Long: fmt.Sprintf(`
=======================================================================
=                        %s gitlab                          =
=======================================================================

Command to interract with Gitlab

`, name.ApplicationName),
	PersistentPreRun: checkGitlabArguments,
}

// Initialize subcommand
func init() {
	rootCmd.AddCommand(gitlabCmd)

	// Common Gitlab flags
	// -> Flags for calling Gitlab API
	gitlabCmd.PersistentFlags().StringVarP(&apiUserToken, "api-token", "t", "", fmt.Sprintf("valid private or personal token to call API methods which require authentication <%s_%s>", config.ENV_PREFIX, "API_TOKEN"))
	gitlabCmd.PersistentFlags().StringVarP(&baseDomain, "domain", "", "gitlab.com", fmt.Sprintf("the domain where gitlab lives <%s_%s>", config.ENV_PREFIX, "DOMAIN"))
	gitlabCmd.PersistentFlags().StringVarP(&gid, "group-id", "g", "", fmt.Sprintf("retrieve all projects under the given group ID <%s_%s>", config.ENV_PREFIX, "GROUP_ID"))
	gitlabCmd.PersistentFlags().StringVarP(&destination, "destination", "", "", fmt.Sprintf("directory destination where projects will be clone <%s_%s>", config.ENV_PREFIX, "DESTINATION"))

	// Define persistent flags
	if err := gitlabCmd.MarkPersistentFlagRequired("api-token"); err != nil {
		logger.Fatalw(err.Error())
	}

	// Add EnvName Param to config
	config.AddEnvParam("API_TOKEN")
	config.AddEnvParam("DOMAIN")
	config.AddEnvParam("GROUP_ID")
	config.AddEnvParam("DESTINATION")
}

// checkGitlabArguments : check CLI args
func checkGitlabArguments(cmd *cobra.Command, args []string) {

	// Check that Gid is a integer
	if gid != "" {
		if _, err := strconv.ParseInt(gid, 10, 32); err != nil {
			logger.Infof("error : '%s' is invalid\n-> group-id must be an integer", gid)
			os.Exit(1)
		}
	}
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
