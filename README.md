# gh-jira-changelog

`gh` cli plugin to generate changelog from jira

## Content
  * [Installation](#installation)
    * [Using gh cli, as extension](#using-gh-cli,-as-extension)
    * [MacOS using `homebrew`](#macos-using-`homebrew`)
    * [Go Toolchain](#go-toolchain)
  * [Verify installation](#verify-installation)
  * [Usage](#usage)
    * [Generating Changelog](#generating-changelog)


### Installation

Feel free to use any of the available methods

#### Using gh cli, as extension
```sh
gh extension install handofgod94/gh-jira-changelog
```

#### MacOS using `homebrew`
```sh
brew install handofgod94/tap/gh-jira-changelog
```

#### Go Toolchain
```sh
go install github.com/handofgod94/gh-jira-changelog@latest
```
The go binary will be installed in `$GOPATH/bin`

### Verify installation

`$ gh-jira-changelog version`
```
dev
```

### Usage

`$ gh-jira-changelog --help`
```
Most of our changelog tools solely focus on commits. While the orgs usually use jira to track issues.
When generating changelog why not combine both commits and jira issues to generate a changelog.

This can also work as a plugin for "gh" cli

Usage:
  gh-jira-changelog [command]

Available Commands:
  auth        Authenticate wit jira
  completion  Generate the autocompletion script for the specified shell
  generate    Generates changelog
  help        Help about any command
  version     Current version of generator

Flags:
  -t, --api_token string   API token for jira
  -u, --base_url string    base url where jira is hosted
      --config string      config file (default is ./.jira_changelog.yaml)
      --email_id string    email id of the user
  -h, --help               help for gh-jira-changelog
  -v, --log_level string   log level. options: debug, info, warn, error (default "error")

Use "gh-jira-changelog [command] --help" for more information about a command.
```

#### Generating Changelog

`$ gh-jira-changelog generate --help`
```
Generates changelog

Usage:
  gh-jira-changelog generate [flags]

Examples:
#using as a standalone tool
gh-jira-changelog generate \
	--base_url="<you-atlassian-url>" \
	--from="<git-ref>" \
	--to="<git-ref>" \
	--api_token="<jira-api-token>" \
	--email_id="jira-email-id"

# using config file
# all the jira config such as (base_url, api_token, email_id) can be stored in a config file
gh-jira-changelog generate --config="<path-to-config-file>.yaml" --from=<git-ref> --to=<git-ref>

# using env variables
# all the jira config such as (base_url, api_token, email_id) can be provided by env variables
BASE_URL=<you-atlassian-url> API_TOKEN=<jira-api-token> gh-jira-changelog generate --from=<git-ref> --to=<git-ref>

# generating changelog between 2 git tags
gh-jira-changelog generate --config="<path-to-config-file>.yaml" --from="v0.1.0" --to="v0.2.0"


# Using it as GH plugin
# assuming jira plugin installed
gh jira-changelog generate --config="<path-to-config-file>.yaml" --from="v0.1.0" --to="v0.2.0"

Flags:
      --from string       Git ref to start from
  -h, --help              help for generate
      --to string         Git ref to end at (default "main")
      --write_to string   File stream to write the changelog (default "/dev/stdout")

Global Flags:
  -t, --api_token string   API token for jira
  -u, --base_url string    base url where jira is hosted
      --config string      config file (default is ./.jira_changelog.yaml)
      --email_id string    email id of the user
  -v, --log_level string   log level. options: debug, info, warn, error (default "error")
```
