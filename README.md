![GitHub](https://img.shields.io/github/license/ptavares/go-git-projects)
[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit&logoColor=white)](https://github.com/pre-commit/pre-commit)
![Release](https://img.shields.io/badge/Release_version-0.0.0-blue)


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
    - [Configuration file](#configuration-file)
  - [License](#license)

<!--TOC-->

## Project Information

This project aims to facilitate the cloning of all Git projects authorized to a user.

## Already availabe

- [X] Gitlab clone
- [ ] GitHub clone

## Usage

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
  group_id: 10              # Id of the group where projects are stored

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
