# washhub - an experimental external Github plugin for [Wash](https://puppetlabs.github.io/wash/)

## TL;DR

`washhub` integrates your Github content into your [Wash](https://puppetlabs.github.io/wash/) filesystem.
You can navigate your organisations, repos, directories and files.

## Installation

1. Of course you already have Wash installed and working on your local system. If not, see [Wash Docs](https://puppetlabs.github.io/wash/) for instructions on how to get up and running!
1. Make sure you have Go installed on your system (tested with go 1.12.5)
1. Clone the `washhub` repo and build the plugin. The executable name needs to be `github`:

    ```bash
    git clone https://github.com/timidri/washhub
    cd washhub
    go build -o github
    ```

1. Configure your github token in `~/.washhub.yaml`:

    ```bash
    cat > ~/.washhub.yaml
    github_token: <my_github_token>
    ```

    Alternatively, you can use user credentials:

    ```bash
    cat > ~/.washhub.yaml
    github_user: <my_github_user>
    github_password: <my_github_password>
    ```

1. Configure Wash to use the washhub plugin by editing the `wash.yaml` configuration file (`~/.puppetlabs/wash/wash.yaml` by default):

    ```yaml
    external-plugins:
      - script: '/path/to/washhub/github'
    ```

1. Enjoy!

## How to use

Washhub supports the following Wash actions:

* `list`: organizations, repos and directories inside repos
* `read`: files

## Known bugs

Listing orgs with lots of repos can result in an empty repository listing. When running wash in debug mode (`wash --loglevel debug`) the following error is displayed:

```quote
DEBU FUSE: Find .git in /github/<org> errored: script returned a non-zero exit code of 1. stderr output: Error: GET https://api.github.com/orgs/<org>/repos?per_page=99: 502 Server Error []
```

This erroneous behaviour is documented [here](https://github.com/google/go-github/issues/999) and I didn't find a work-around for now. Ideas welcome!

## Limitations of `washhub`

* Supported are repos, directories and files, other objects such as projects, issues or commits are not supported
* No `metadata` support yet
* `stream` or `exec` not supported because Github doesn't support that
* non-authenticated operation not supported
* You can hit Github API limits if you use `washhub` _a lot_
  
## Author

dimitri@puppet.com
