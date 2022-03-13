package exception

import (
	"fmt"
)

// ParseConfigFileError : custom exception
func ParseConfigFileError(configFileName string, err error) error {
	return fmt.Errorf("unable to parse configuration file [%s] : [%w]", configFileName, err)
}

// InitGitlabClientError : custom exception
func InitGitlabClientError(err error) error {
	return fmt.Errorf("unable to configure a Gitlab client : [%w]", err)
}

// APICallError : custom exception
func APICallError(message string) error {
	return fmt.Errorf("%s", message)
}

// APICallError : custom exception
func SimpleError(message string, err error) error {
	return fmt.Errorf("%s : %w", message, err)
}

// FetchingURLError : custom exception
func FetchingURLError(url string, err error) error {
	return fmt.Errorf("error fetching URL %s : %w", url, err)
}

// ReadPrivateKeyFileError : custom exception
func ReadPrivateKeyFileError(filepath string, err error) error {
	return fmt.Errorf("unable to load ssh private key file from filepath [%s] : %w", filepath, err)
}

// GeneratePublicKeysError : custom exception
func GeneratePublicKeysError(err error) error {
	return fmt.Errorf("unable to generate ssh public keys : %w", err)
}

// GetHomeDirError : custom exception
func GetPWDDirError(err error) error {
	return fmt.Errorf("unable to get current directory : %w", err)
}

// CreateDirError : custom exception
func CreateDirError(directoryPath string, err error) error {
	return fmt.Errorf("unable to get create directory %s : %w", directoryPath, err)
}
