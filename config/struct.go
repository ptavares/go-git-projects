package config

import (
	"fmt"
	"git-projects/helper"
	"strings"
)

// Config: represents a git-projects configuration
type Config struct {
	ApiToken    string      `mapstructure:"api_token"`
	Domain      string      `mapstructure:"domain,omitempty"`
	CloneConfig CloneConfig `mapstructure:"clone_config,omitempty"`
	BasicAuth   BasicAuth   `mapstructure:"clone_basic_auth,omitempty"`
	SSHAuth     SSHAuth     `mapstructure:"clone_ssh_auth,omitempty"`
}

// CloneConfig: contains clone custom configuration
type CloneConfig struct {
	Destination string `mapstructure:"destination,omitempty"`
	GroupID     string `mapstructure:"group_id,omitempty"`
}

// BasicAuth: define clone basic authentication
type BasicAuth struct {
	UserName string `mapstructure:"user_name,omitempty"`
	Password string `mapstructure:"password,omitempty"`
	Token    string `mapstructure:"user_token,omitempty"`
}

// SSHAuth: define clone ssh authentication
type SSHAuth struct {
	PrivateKeyPath     string `mapstructure:"key_path"`
	PrivateKeyPassword string `mapstructure:"key_password,omitempty"`
}

// String for Config (mask sensitive values)
func (c *Config) String() string {
	s := new(strings.Builder)
	s.WriteString("{ ")
	s.WriteString("ApiToken: \"")
	s.WriteString(helper.MaskPassword(c.ApiToken))
	s.WriteString("\" Domain: \"")
	s.WriteString(c.Domain)
	s.WriteString("\" CloneConfig: ")
	s.WriteString(fmt.Sprintf("%+v", c.CloneConfig))
	s.WriteString("\" BasicAuth: ")
	s.WriteString(fmt.Sprintf("%+v", c.BasicAuth))
	s.WriteString("\" SSHAuth: ")
	s.WriteString(fmt.Sprintf("%+v", c.SSHAuth))
	return s.String()
}
