package cmd

import (
	"fmt"
	"git-projects/config"
	"git-projects/internal/name"

	"github.com/spf13/cobra"
)

// githubCmd represents the config command
var githubCmd = &cobra.Command{
	Use:   "github",
	Short: fmt.Sprintf("Perform %s Github actions", name.ApplicationName),
	Long: fmt.Sprintf(`
=======================================================================
=                        %s github                          =
=======================================================================

Command to interract with Github

`, name.ApplicationName),
	PersistentPreRunE: checkGithubArguments,
}

const default_github_domain = "github.com"

// Initialize subcommand
func init() {
	rootCmd.AddCommand(githubCmd)

	// Common Github flags
	// -> Flags for calling Github API
	githubCmd.PersistentFlags().StringVarP(&apiUserToken, "api-token", "t", "", fmt.Sprintf("valid private or personal token to call API methods which require authentication <%s_%s>", config.ENV_PREFIX, "API_TOKEN"))
	githubCmd.PersistentFlags().StringVarP(&baseDomain, "domain", "", default_github_domain, fmt.Sprintf("the domain where github lives <%s_%s>", config.ENV_PREFIX, "DOMAIN"))
	githubCmd.PersistentFlags().StringVarP(&gid, "group-id", "g", "", fmt.Sprintf("ID of the group who's repos should be cloned <%s_%s>", config.ENV_PREFIX, "GROUP_ID"))
	githubCmd.PersistentFlags().StringVarP(&destination, "destination", "", "", fmt.Sprintf("directory destination where projects will be clone, default is current directory <%s_%s>", config.ENV_PREFIX, "DESTINATION"))

	// Define persistent flags
	if err := githubCmd.MarkPersistentFlagRequired("api-token"); err != nil {
		logger.Fatalw(err.Error())
	}

	// Add EnvName Param to config
	config.AddEnvParam("API_TOKEN")
	config.AddEnvParam("DOMAIN")
	config.AddEnvParam("GROUP_ID")
	config.AddEnvParam("DESTINATION")
}

// checkGithubArguments : check CLI args
func checkGithubArguments(cmd *cobra.Command, args []string) error {
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
	fmt.Println(localDomain, " ", baseDomain)
	if localDomain != "" {
		baseDomain = localDomain
	} else {
		baseDomain = default_github_domain
	}

	return nil
}
