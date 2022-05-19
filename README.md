![GitHub](https://img.shields.io/github/license/ptavares/go-git-projects)
[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit&logoColor=white)](https://github.com/pre-commit/pre-commit)
![Release](https://img.shields.io/badge/Release_version-0.2.1-blue)


# go-git-projects

Project to clone Github/Gitlab projects in pure Golang language

## Table of content

_This documentation section is generated automatically_

<!--TOC-->

- [go-git-projects](#go-git-projects)
  - [Table of content](#table-of-content)
  - [Project Information](#project-information)
  - [Already availabe](#already-availabe)
  - [Usage](#usage)
    - [Available environment vars](#available-environment-vars)
    - [CLI](#cli)
    - [Configuration file](#configuration-file)
  - [License](#license)

<!--TOC-->

## Project Information

This project aims to facilitate the cloning of all Git projects authorized to a user.

## Already availabe

- [X] Gitlab clone
- [X] GitHub clone

## Usage

You can find here a list of common usage for this application

### Available environment vars

You can use some environments variables with this CLI to avoid passing arguments.
You don't need to setup all of them, it depends on witch command you want to use.

There are indicated in help command usage, and started with `GIT_PROJECT_` prefix.

### CLI

#### Root Command

```
=======================================================================
=                           git-projects                              =
=======================================================================

A CLI to easyli clone/sync git projects from :
  -> Gitlab
  -> Github

Usage:
  git-projects [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  github      Perform git-projects Github actions
  gitlab      Perform git-projects Gitlab actions
  help        Help about any command
  version     Show the git-projects version information

Flags:
  -c, --config string   config yaml file (default are : ${HOME}/.git-projects[.yaml] or ${PWD}/.git-projects[.yaml])
  -d, --debug           show debug message
  -h, --help            help for git-projects

Use "git-projects [command] --help" for more information about a command.
```

#### Github SubCommand

```
=======================================================================
=                        git-projects github                          =
=======================================================================

Command to interract with Github

Usage:
  git-projects github [command]

Available Commands:
  clone       Perform git-projects Github clone actions

Flags:
  -t, --api-token string      valid private or personal token to call API methods which require authentication <GIT_PROJECTS_API_TOKEN>
      --destination string    directory destination where projects will be clone, default is current directory <GIT_PROJECTS_DESTINATION>
      --domain string         the domain where github lives <GIT_PROJECTS_DOMAIN> (default "github.com")
  -h, --help                  help for github
  -o, --organization string   organization name who's repos should be cloned <GIT_PROJECTS_ORGANIZATION>
  -u, --user string           user name who's repos should be cloned <GIT_PROJECTS_USER>

Global Flags:
  -c, --config string   config yaml file (default are : ${HOME}/.git-projects[.yaml] or ${PWD}/.git-projects[.yaml])
  -d, --debug           show debug message

Use "git-projects github [command] --help" for more information about a command.
```

##### Github Clone SubCommand

1. Clone from HTTP

```
=======================================================================
=                   git-projects github clone HTTP                    =
=======================================================================

Command to clone Github projects from repositories's HTTP URL

Usage:
  git-projects github clone http [flags]

Examples:

  # =========================
  # Cloning all user projects
  # =========================
  -> Using same token for cloning
  git-projects github clone http -t <token> -u <github_username>
  git-projects github clone http -t <token> -u <github_username> --destination /tmp/dest

  -> Using basic auth token
  git-projects github clone http -t <token> -u <github_username> --basic-auth-token <cloning_token>

  -> Using basic auth user/pwd
  git-projects github clone http -t <token> -u <github_username> --basic-auth-username <cloning_user_name> --basic-auth-password <cloning_user_password>

  # =================================
  # Cloning all organization projects
  # =================================
  -> Using same token for cloning
  git-projects github clone http -t <token> -o <github_organization>
  git-projects github clone http -t <token> -o <github_organization> --destination ./github

  -> Using basic auth token
  git-projects github clone http -t <token> -o <github_organization> --basic-auth-token <cloning_token>

  -> Using basic auth user/pwd
  git-projects github clone http -t <token> -o <github_organization> --basic-auth-username <cloning_user_name> --basic-auth-password <cloning_user_password>

  # ======================================
  # Cloning all projects using config file
  # ======================================
  -> Config in  default location
  git-projects github clone http

  -> Specify custom file
  git-projects github clone http -c <path_to_config_file>

  # ======================================
  # Entreprise Github domain
  # ======================================
  -> Use "--domain" flag
  git-projects github clone http -d <entreprise_domain>


Flags:
      --basic-auth-password string   password related to 'basic-auth-username' <GIT_PROJECTS_BASIC_AUTH_PWD>
      --basic-auth-token string      token to use to clone repository throw HTTP URL if different from 'api-token' <GIT_PROJECTS_BASIC_AUTH_TOKEN>
      --basic-auth-username string   username to use to clone repository throw HTTP URL <GIT_PROJECTS_BASIC_AUTH_USR>
  -h, --help                         help for http

Global Flags:
  -t, --api-token string      valid private or personal token to call API methods which require authentication <GIT_PROJECTS_API_TOKEN>
  -c, --config string         config yaml file (default are : ${HOME}/.git-projects[.yaml] or ${PWD}/.git-projects[.yaml])
  -d, --debug                 show debug message
      --destination string    directory destination where projects will be clone, default is current directory <GIT_PROJECTS_DESTINATION>
      --domain string         the domain where github lives <GIT_PROJECTS_DOMAIN> (default "github.com")
  -o, --organization string   organization name who's repos should be cloned <GIT_PROJECTS_ORGANIZATION>
  -u, --user string           user name who's repos should be cloned <GIT_PROJECTS_USER>
```

2. Clone form SSH

```
=======================================================================
=                   git-projects github clone SSH                   =
=======================================================================

Command to clone Github projects from repositories's SSH URL

Usage:
  git-projects github clone ssh [flags]

Examples:

  # =========================
  # Cloning all user projects
  # =========================
  -> Using same token for cloning
  git-projects github clone ssh -t <token> -u <github_username>
  git-projects github clone ssh -t <token> -u <github_username> --destination /tmp/dest

  -> Using ssh auth file
  git-projects github clone ssh -t <token> -u <github_username> --ssh-private-key-path <path_to_privat_key> --ssh-private-key-password <optional_key_password>

  # =================================
  # Cloning all organization projects
  # =================================
  -> Using same token for cloning
  git-projects github clone ssh -t <token> -o <github_organization>
  git-projects github clone ssh -t <token> -o <github_organization> --destination ./github

  -> Using ssh auth file
  git-projects github clone ssh -t <token> -o <github_organization> --ssh-private-key-path <path_to_privat_key> --ssh-private-key-password <optional_key_password>

  # ======================================
  # Cloning all projects using config file
  # ======================================
  -> Config in  default location
  git-projects github clone ssh

  -> Specify custom file
  git-projects github clone ssh -c <path_to_config_file>

  # ======================================
  # Entreprise Github domain
  # ======================================
  -> Use "--domain" flag
  git-projects github clone ssh -d <entreprise_domain>


Flags:
  -h, --help                              help for ssh
      --ssh-private-key-password string   optional password to decrypt private key <GIT_PROJECTS_SSH_KEY_PWD>
      --ssh-private-key-path string       path to private key file used to clone repository throw SSH URL <GIT_PROJECTS_SSH_KEY_PATH>

Global Flags:
  -t, --api-token string      valid private or personal token to call API methods which require authentication <GIT_PROJECTS_API_TOKEN>
  -c, --config string         config yaml file (default are : ${HOME}/.git-projects[.yaml] or ${PWD}/.git-projects[.yaml])
  -d, --debug                 show debug message
      --destination string    directory destination where projects will be clone, default is current directory <GIT_PROJECTS_DESTINATION>
      --domain string         the domain where github lives <GIT_PROJECTS_DOMAIN> (default "github.com")
  -o, --organization string   organization name who's repos should be cloned <GIT_PROJECTS_ORGANIZATION>
  -u, --user string           user name who's repos should be cloned <GIT_PROJECTS_USER>
```

#### Gitlab SubCommand

```
=======================================================================
=                        git-projects gitlab                          =
=======================================================================

Command to interract with Gitlab

Usage:
  git-projects gitlab [command]

Available Commands:
  clone       Perform git-projects Gitlab clone actions

Flags:
  -t, --api-token string     valid private or personal token to call API methods which require authentication <GIT_PROJECTS_API_TOKEN>
      --destination string   directory destination where projects will be clone, default is current directory <GIT_PROJECTS_DESTINATION>
      --domain string        the domain where gitlab lives <GIT_PROJECTS_DOMAIN> (default "gitlab.com")
  -g, --group-id string      ID of the group who's repos should be cloned <GIT_PROJECTS_GROUP_ID>
  -h, --help                 help for gitlab

Global Flags:
  -c, --config string   config yaml file (default are : ${HOME}/.git-projects[.yaml] or ${PWD}/.git-projects[.yaml])
  -d, --debug           show debug message

Use "git-projects gitlab [command] --help" for more information about a command.
```


##### Github Clone SubCommand

1. Clone from HTTP

```
======================================================================
=                   git-projects gitlab clone HTTP                    =
=======================================================================

Command to clone Gitlab projects from repositories's HTTP URL

Usage:
  git-projects gitlab clone http [flags]

Examples:

  # =========================
  # Cloning all user projects
  # =========================
  -> Using same token for cloning
  git-projects gitlab clone http -t <token>
  git-projects gitlab clone http -t <token> --destination /tmp/dest

  -> Using basic auth token
  git-projects gitlab clone http -t <token> --basic-auth-token <cloning_token>

  -> Using basic auth user/pwd
  git-projects gitlab clone http -t <token> --basic-auth-username <cloning_user_name> --basic-auth-password <cloning_user_password>

  # =================================
  # Cloning all GroupId projects
  # =================================
  -> Using same token for cloning
  git-projects gitlab clone http -t <token> -g <gitlab_group_id>
  git-projects gitlab clone http -t <token> -g <gitlab_group_id> --destination ./gitlab

  -> Using basic auth token
  git-projects gitlab clone http -t <token> -g <gitlab_group_id> --basic-auth-token <cloning_token>

  -> Using basic auth user/pwd
  git-projects gitlab clone http -t <token> -g <gitlab_group_id> --basic-auth-username <cloning_user_name> --basic-auth-password <cloning_user_password>

  # ======================================
  # Cloning all projects using config file
  # ======================================
  -> Config in  default location
  git-projects gitlab clone http

  -> Specify custom file
  git-projects gitlab clone http -c <path_to_config_file>

  # ======================================
  # Entreprise Gitlab domain
  # ======================================
  -> Use "--domain" flag
  git-projects gitlab clone http -d <entreprise_domain>


Flags:
      --basic-auth-password string   password related to 'basic-auth-username' <GIT_PROJECTS_BASIC_AUTH_PWD>
      --basic-auth-token string      token to use to clone repository throw HTTP URL if different from 'api-token' <GIT_PROJECTS_BASIC_AUTH_TOKEN>
      --basic-auth-username string   username to use to clone repository throw HTTP URL <GIT_PROJECTS_BASIC_AUTH_USR>
  -h, --help                         help for http

Global Flags:
  -t, --api-token string     valid private or personal token to call API methods which require authentication <GIT_PROJECTS_API_TOKEN>
  -c, --config string        config yaml file (default are : ${HOME}/.git-projects[.yaml] or ${PWD}/.git-projects[.yaml])
  -d, --debug                show debug message
      --destination string   directory destination where projects will be clone, default is current directory <GIT_PROJECTS_DESTINATION>
      --domain string        the domain where gitlab lives <GIT_PROJECTS_DOMAIN> (default "gitlab.com")
  -g, --group-id string      ID of the group who's repos should be cloned <GIT_PROJECTS_GROUP_ID>
```

2. Clone from SSH

```
=======================================================================
=                   git-projects gitlab clone SSH                   =
=======================================================================

Command to clone Gitlab projects from repositories's SSH URL

Usage:
  git-projects gitlab clone ssh [flags]

Examples:

  # =========================
  # Cloning all user projects
  # =========================
  -> Using same token for cloning
  git-projects gitlab clone ssh -t <token>
  git-projects gitlab clone ssh -t <token> --destination /tmp/dest

  -> Using ssh auth file
  git-projects gitlab clone ssh -t <token> --ssh-private-key-path <path_to_privat_key> --ssh-private-key-password <optional_key_password>

  # =================================
  # Cloning all GroupId projects
  # =================================
  -> Using same token for cloning
  git-projects gitlab clone ssh -t <token> -g <gitlab_group_id>
  git-projects gitlab clone ssh -t <token> -g <gitlab_group_id> --destination ./gitlab

  -> Using ssh auth file
  git-projects gitlab clone ssh -t <token> -g <gitlab_group_id> --ssh-private-key-path <path_to_privat_key> --ssh-private-key-password <optional_key_password>

  # ======================================
  # Cloning all projects using config file
  # ======================================
  -> Config in  default location
  git-projects gitlab clone ssh

  -> Specify custom file
  git-projects gitlab clone ssh -c <path_to_config_file>

  # ======================================
  # Entreprise Github domain
  # ======================================
  -> Use "--domain" flag
  git-projects gitlab clone ssh -d <entreprise_domain>


Flags:
  -h, --help                              help for ssh
      --ssh-private-key-password string   optional password to decrypt private key <GIT_PROJECTS_SSH_KEY_PWD>
      --ssh-private-key-path string       path to private key file used to clone repository throw SSH URL <GIT_PROJECTS_SSH_KEY_PATH>

Global Flags:
  -t, --api-token string     valid private or personal token to call API methods which require authentication <GIT_PROJECTS_API_TOKEN>
  -c, --config string        config yaml file (default are : ${HOME}/.git-projects[.yaml] or ${PWD}/.git-projects[.yaml])
  -d, --debug                show debug message
      --destination string   directory destination where projects will be clone, default is current directory <GIT_PROJECTS_DESTINATION>
      --domain string        the domain where gitlab lives <GIT_PROJECTS_DOMAIN> (default "gitlab.com")
  -g, --group-id string      ID of the group who's repos should be cloned <GIT_PROJECTS_GROUP_ID>
```

### Configuration file

#### Location

By default, this application will check for a configuration file, named **.git-projects[.yaml]** in two directories :

1. User homedir : `$HOME/.git-projects[.yaml]`
2. Current directory (from where the application is run) : `$PWD/.git-projects[.yaml]`

-> You can specify the location of your configuration file using the flag **`--config`** or the shortest way **`-c`**

#### Structure

The configuration file must have the below structure :

```yaml
api_token: xxxxxxxxxx       # Token api to use to call Git[Hub/Lab] API
domain: my.custom.git.com   # The custom domain to use if different from default (gitlab.com or github.com)
# Custom clone configuration
clone_config:
  destination: /tmp/test    # Directory where all projects will be cloned
  # Gitlab specific configuration
  group_id: 10              # Id of the group where projects are stored
  # Github specific configuration (one of user or organization)
  user: ptavares            # github user where projects are stored
  organization: myorg       # github organization where projects are stored

# If the api_token is not the one to use for cloning projects,
# you can specify basic or ssh authentication method

# Basic authentication
# You can configure a user/pws tuple or a user_token
clone_basic_auth:
  user_name: user           # Username for basic authentication / password must be filled
  password: pwd             # Password for basic auyhentication when using user_name
  user_token: xxxxxx        # User token to use if don't user user_name/password tuple

# SSH authentication
clone_ssh_auth:
  key_path: /home/$USER/.ssh/id_rsa  # Path to the private key
  key_password: xxxxxx               # Password to read private key if private key is password protected
```

## License

[MIT](LICENCE)
