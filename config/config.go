package config

import (
	"git-projects/helper"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	// Prefix for application environment variables
	ENV_PREFIX string = "GIT_PROJECTS"

	// envVarNames is a list that contains all supported Environement Variables to pass to CLI
	envVarNames []string
)

// AddEnvParam : Add an Environment Param
func AddEnvParam(envVarName string) {
	envVarNames = append(envVarNames, envVarName)
}

// LoadConfig : Load Configuration using viper (file + environment).
// You can use GetConfig next
func LoadConfig() error {
	logger := helper.GetLogger()

	logger.Debug("start LoadConfig()")

	viper.AutomaticEnv() // read in environment variables that match
	// var starting with BPI
	viper.SetEnvPrefix(ENV_PREFIX)

	logger.With(
		zap.String("EnvPrefix", ENV_PREFIX),
		zap.Strings("envVarNames", envVarNames),
	).Debug("Binding Environment variables")

	// List of environment Name
	for _, envVarName := range envVarNames {
		helper.HandleError(viper.BindEnv(envVarName))
	}

	logger.Debug("end - LoadConfig() ")
	return nil
}
