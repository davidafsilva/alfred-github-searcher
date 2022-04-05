# alfred-github-searcher

A simple [Alfred workflow](https://www.alfredapp.com/help/workflows/) that
enables searching through the list of repositories that you've contributed to as
well as pull requests created or pending your review.

## Installation

> :warning: Requirement
>
> In order to install this workflow you'll need a valid license for
> [Alfred's powerpack](https://www.alfredapp.com/powerpack/)

1. Download the latest version of the workflow from the
   [releases page](https://github.com/davidafsilva/alfred-github-searcher/releases)
2. Install the `.alfredworkflow` file
3. Create a [personal access token](https://github.com/settings/tokens) with
   the `repo` scope enabled
    1. Be sure to authorize the token on your organizations in order to enable
       searching on its private repositories
4. Set up the environment variable named `ags_github_token` on the workflow
   settings with a valid GitHub token as its value

## Usage / Available Commands

| Command    | Arguments  | Description                                     |
|------------|------------|-------------------------------------------------|
| `repo`     | `[filter]` | Searches for repositories you've contributed to |
| `reposync` |            | Forces the repository synchronization           |
| `prc`      | `[filter]` | Searches for PRs that you've created            |
| `prr`      | `[filter]` | Searches for PRs that are pending your review   |
| `pr`       | `[filter]` | Both of the above                               |
| `prsync`   |            | Forces the pull request synchronization         |
| `ghu`      |            | Checks for and prompts for an workflow update   |

## Configuration

| Name                                | Default | Description                                                                 |
|-------------------------------------|---------|-----------------------------------------------------------------------------|
| `ags_github_token`                  |         | Your personal GitHub Token                                                  |
| `ags_prs_refresh_interval`          | `30m`   | Refresh interval for the local pull request data                            |
| `ags_repositories_refresh_interval` | `120h`  | Refresh interval for the local repository data                              |
| `ags_show_owner_image`              | `true`  | Whether or not to show the owner's image next to the repository suggestions |
| `ags_theme`                         | `light` | Choose between `light` icons or `dark` icons                                |
