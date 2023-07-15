# gh-jira-changelog

`gh` cli plugin to generate changelog from jira

## Content
  * [Installation](#installation)
  * [Verify installation](#verify-installation)
  * [Usage](#usage)
    * [Generating Changelog](#generating-changelog)


### Installation

WIP

### Verify installation

`$ gh-jira-changelog version`
```
v0.1.0
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
      --config string   config file (default is ./.jira_changelog.yaml)
  -h, --help            help for gh-jira-changelog

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
	--project_name="<jira-project-name>" \
	--from="<git-ref>" \
	--to="<git-ref>" \
	--api_token="<jira-api-token>" \
	--email_id="jira-email-id"

# using config file
# all the jira config such as (base_url, project_name, api_token, email_id) can be stored in a config file
gh-jira-changelog generate --config="<path-to-config-file>.yaml" --from=<git-ref> --to=<git-ref>

# using env variables
# all the jira config such as (base_url, project_name, api_token, email_id) can be provided by env variables
BASE_URL=<you-atlassian-url> PROJECT_NAME=<jira-project-name> API_TOKEN=<jira-api-token> gh-jira-changelog generate --from=<git-ref> --to=<git-ref>

# generating changelog between 2 git tags
gh-jira-changelog generate --config="<path-to-config-file>.yaml" --from="v0.1.0" --to="v0.2.0"


# Using it as GH plugin
# assuming jira plugin installed
gh jira-changelog generate --config="<path-to-config-file>.yaml" --from="v0.1.0" --to="v0.2.0"

Flags:
  -t, --api_token string      API token for jira
  -u, --base_url string       base url where jira is hosted
      --email_id string       email id of the user
      --from string           Git ref to start from
  -h, --help                  help for generate
  -v, --log_level string      log level. options: debug, info, warn, error (default "error")
  -p, --project_name string   Project name in jira. usually the acronym
      --to string             Git ref to end at (default "main")
      --write_to string       File stream to write the changelog (default "/dev/stdout")

Global Flags:
      --config string   config file (default is ./.jira_changelog.yaml)
```
