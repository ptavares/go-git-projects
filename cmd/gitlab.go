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
	PersistentPreRunE: checkGitlabArguments,
}

// Initialize subcommand
func init() {
	rootCmd.AddCommand(gitlabCmd)

	// Common Gitlab flags
	// -> Flags for calling Gitlab API
	gitlabCmd.PersistentFlags().StringVarP(&apiUserToken, "api-token", "t", "", fmt.Sprintf("valid private or personal token to call API methods which require authentication <%s_%s>", config.ENV_PREFIX, "API_TOKEN"))
	gitlabCmd.PersistentFlags().StringVarP(&baseDomain, "domain", "", "gitlab.com", fmt.Sprintf("the domain where gitlab lives <%s_%s>", config.ENV_PREFIX, "DOMAIN"))
	gitlabCmd.PersistentFlags().StringVarP(&gid, "group-id", "g", "", fmt.Sprintf("ID of the group who's repos should be cloned <%s_%s>", config.ENV_PREFIX, "GROUP_ID"))
	gitlabCmd.PersistentFlags().StringVarP(&destination, "destination", "", "", fmt.Sprintf("directory destination where projects will be clone, default is current directory <%s_%s>", config.ENV_PREFIX, "DESTINATION"))

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
func checkGitlabArguments(cmd *cobra.Command, args []string) error {
	if err := rootCmd.PersistentPreRunE(cmd, args); err != nil {
		return err
	}

	fillStringParam("api_token", config.GetConfig().ApiToken, &apiUserToken)
	// Remove persistent flag for api-token if it's present
	if apiUserToken != "" {
		cmd.Flag("api-token").Changed = true
	}

	fillStringParam("destination", config.GetConfig().CloneConfig.Destination, &destination)
	fillStringParam("group_id", config.GetConfig().CloneConfig.GroupID, &gid)

	var localDomain = ""
	fillStringParam("domain", config.GetConfig().Domain, &localDomain)
	if localDomain != "" {
		baseDomain = localDomain
	}

	// Check that Gid is a integer
	if gid != "" {
		if _, err := strconv.ParseInt(gid, 10, 32); err != nil {
			logger.Infof("error : '%s' is invalid\n-> group-id must be an integer", gid)
			os.Exit(1)
		}
	}

	return nil
}
