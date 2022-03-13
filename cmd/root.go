package cmd

import (
	"fmt"
	"git-projects/config"
	"git-projects/helper"
	"git-projects/internal/name"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Common variable (used by all commands)
var (
	cfgFile           string
	debug             bool
	logger            *zap.SugaredLogger
	apiUserToken      string
	baseDomain        string
	sshPrivateKeyPath string
	sshPrivateKeyPwd  string
	basicAuthUsername string
	basicAuthPassword string
	basicAuthToken    string
	gid               string
	destination       string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:               name.ApplicationName,
	Short:             "A CLI to easyli clone/sync git project",
	PersistentPreRunE: runInit,
	Long: fmt.Sprintf(`
=======================================================================
=                           %s                              =
=======================================================================

A CLI to easyli clone/sync git projects from :
  -> Gitlab
  -> Github

	`, name.ApplicationName),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Define root global command flags
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config yaml file (default are : ${HOME}/.git-projects[.yaml] or ${PWD}/.git-projects[.yaml])")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "show debug message")

	// Run init
	cobra.OnInitialize(initRoot)
}

// initRoot : will load common configuration (logger)
func initRoot() {
	logger = helper.InitLogger(debug)
}

// Run init on root before other pre-run : will load configuration
func runInit(cmd *cobra.Command, args []string) error {
	logger.Debug("start - runInit")

	logger.Debug("Loading configuration...")
	if _, err := config.LoadConfig(cfgFile); err != nil {
		return err
	}
	logger.Debug("end - runInit")

	return nil
}

// Fill paramValue with defaultValue if paramName is empty
func fillStringParam(paramName string, defaultValue string, paramValue *string) {
	if *paramValue == "" {
		*paramValue = helper.Ternary(viper.GetString(paramName) == "", defaultValue, viper.GetString(paramName)).(string)
	}
}
