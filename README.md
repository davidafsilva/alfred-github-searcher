# alfred-github-repositories

A simple [Alfred workflow](https://www.alfredapp.com/help/workflows/) that
enables searching through the list of repositories that you've contributed to.

<img alt="preview" src="preview.png"/>

## Installation

> :warning: Requirement
>
> In order to install this workflow you'll need a valid license for
> [Alfred's powerpack](https://www.alfredapp.com/powerpack/)

1. Download the latest version of the workflow from the
   [releases page](https://github.com/davidafsilva/alfred-github-repositories/releases)
2. Install the `.alfredworkflow` file
3. Create a [personal access token](https://github.com/settings/tokens)
   with the `repo` scope enabled
    1. Be sure to authorize the token on your organizations in order to enable
       searching on its private repositories
4. Add an environment variable named `alfred_github_repos_token` on the workflow
   settings with the GitHub token as its value
5. Synchronize the remote repositories via the `ghs` command (keyword)
6. Search through the repositories with `gh <repo>`
