package config

import (
	"fmt"
	"git-projects/exception"
	"git-projects/helper"
	"git-projects/internal/name"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	// Prefix for application environment variables
	ENV_PREFIX string = "GIT_PROJECTS"

	// envVarNames is a list that contains all supported Environement Variables to pass to CLI
	envVarNames []string

	// config is a Struct that contains the loaded configuration
	config *Config
)

// AddEnvParam : Add an Environment Param
func AddEnvParam(envVarName string) {
	envVarNames = append(envVarNames, envVarName)
}

// GetConfig : Return Configuration
func GetConfig() *Config {
	return config
}

// LoadConfig : Load Configuration using viper (file + environment).
// You can use GetConfig next
func LoadConfig(configFileName string) (*Config, error) {
	logger := helper.GetLogger()

	logger.With(
		zap.String("configFileName", configFileName),
	).Debug("start - LoadConfig(configFileName string)")

	// Init empty config
	config = &Config{}

	if configFileName != "" {
		logger.With(
			zap.String("configFileName", configFileName),
		).Debug("Will use config file from the flag (cli)")

		viper.SetConfigFile(configFileName)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		// Find curent dir
		currentDir, err := os.Getwd()
		cobra.CheckErr(err)

		logger.With(
			zap.String("userHomeDir", home),
			zap.String("currentDir", currentDir),
		).Debug("Configfile empty, try to find it in default directories")

		// Search config in home directory and current.
		viper.AddConfigPath(home)
		viper.AddConfigPath(currentDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName(fmt.Sprintf(".%s", name.ApplicationName))
	}

	viper.AutomaticEnv() // read in environment variables that match
	// var starting with ENV_PREFIX
	viper.SetEnvPrefix(ENV_PREFIX)

	logger.With(
		zap.String("EnvPrefix", ENV_PREFIX),
		zap.Strings("envVarNames", envVarNames),
	).Debug("Binding Environment variables")

	// List of environment Name
	for _, envVarName := range envVarNames {
		helper.HandleError(viper.BindEnv(envVarName))
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		yellow := color.New(color.Bold, color.FgYellow).SprintFunc()
		bold := color.New(color.Bold).SprintFunc()
		logger.Info(yellow("------------------------"))
		logger.Infof("%s : %s", bold("Using configuration file"), viper.ConfigFileUsed())
		logger.Info(yellow("------------------------"))

		err = viper.Unmarshal(&config)
		if err != nil {
			return config, exception.ParseConfigFileError(viper.ConfigFileUsed(), err)
		}

		logger.With(
			zap.String("config", config.String()),
		).Debug("Loaded Configuration")
	}

	logger.With(
		zap.Strings("allkeys", viper.AllKeys()),
	).Debug("All configuration keys, including environment variables")

	logger.Debug("end - LoadConfig(configFileName string) ")

	return config, nil
}
