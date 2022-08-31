# Kojirō

It is a [Klio](https://github.com/g2a-com/klio) compliant application for creating jira releases based on git commit history.
Based on a given tag (must be semver) kojiro looks for a previous tag in commit history.
It browses all the commit messages in-between trying to find jira issues (e.g. *ABC-123*).
Then, it takes all those issues and puts them in a new release named after the given version.

Idea and application is based (and forked from) [jira-versioner](https://github.com/psmarcin/jira-versioner) 

## Getting Started

> IMPORTANT!
> Kojiro works only with Jira Server instances. It does not work with Jira Cloud.

### Prerequisites

Things that you need to have before we start:

* Jira project code that a release will be created in
* Jira user (as an email) with rights to write to Jira project
* A password for given user (only for application users)
* Git repository
* At least two Git tags in rage

### Installing



## Usage

To check current usage of the command simply type:
```console
klio kojiro --help
```

### Example

```console
klio kojiro -e jira@example.com -k SOME_TOKEN -p 10003 -v v1.1.0 -t v1.1.0 -u https://example.atlassian.net
```

### Explained more in depth

Here is our git log history:

```console
05e5705322cc2d9daf7fb376a8c5e9cbd039b257 (HEAD -> master, tag: v2.1.0, origin/master) chore: remove unnecessary string conversion
9bf13576317845cd7d10980d62afe719872ceb01 feat: error logs contains command output JR-4
831e4c253829dbc12683baa5b4d494aa3524f39f feat: jira version not required, default tag JR-13
533569497a68f04674e43e23d06fb9c1f0b3b958 docs: update readme.md with new name JR-2
1e1dd3131aeed3611e70d5f329989c1a09371822 (tag: v2.0.1) chore: rename to jira-versioner JR-3
aeae65755553d03920d7cd7c4a5fdb40a02d7c57 docs: update command name JR-3
e874a9c6162fd102b9de926397a855c1b0dbd880 docs: README.md file JR-2
d15916037b0a6ca04776e474ac461e767631c838 (tag: v2.0.0) feat: consistent arguments name
cb59ea7f0bc3efb8b92de87cd88b589024d18ee7 (tag: v1.1.0) feat: JR-40 argument to run git commands in different path
2e6d61dee0c4ed3a0f7f887973dbc326a487675b (tag: v1.0.0) feat: github JR-1 release action
```

> For simplification jira configs are omitted in the following examples 

1. `klio kojiro -t v2.1.0`
    
    Found commits:
    1. `05e5705322cc2d9daf7fb376a8c5e9cbd039b257`
    2. `9bf13576317845cd7d10980d62afe719872ceb01`
    3. `831e4c253829dbc12683baa5b4d494aa3524f39f`
    4. `533569497a68f04674e43e23d06fb9c1f0b3b958`
    
    Found tasks:
    1. `JR-4`
    2. `JR-13`
    3. `JR-2`
    
    Results: 
    1. New Jira release created, if it doesn't exist already - `v2.1.0`.
    2. If a task have been found in a commit messages, and it exists in Jira it is matched with the release.
    
2. `klio kojiro -t v2.0.1`
    
    Found commits:
    1. `1e1dd3131aeed3611e70d5f329989c1a09371822`
    2. `aeae65755553d03920d7cd7c4a5fdb40a02d7c57`
    3. `e874a9c6162fd102b9de926397a855c1b0dbd880`
    
    Found tasks:
    1. `JR-3`
    2. `JR-2`
    
    Results: 
    1. New Jira release is created, if it doesn't exist already - `v2.0.1`.
    2. If a task have been found in a commit messages, and it exists in Jira it is matched with the release.
    3. Each task is matched exactly once with a release despite multiple occurrences in commit messages.
    
3. `klio kojiro -t v2.0.0`
    
    Found commits:
    1. `d15916037b0a6ca04776e474ac461e767631c838`
    
    Found tasks:
    *none*
    
    Results: 
    1. New Jira release is created, if it doesn't exist already - `v2.0.0`.
    2. No tasks found means no links, means empty version is creates.
    
4. `klio kojiro -t v1.1.0 -v 1000.000.000`
    
    Found commits:
    1. `cb59ea7f0bc3efb8b92de87cd88b589024d18ee7`
    
    Found tasks:
    1. `JR-40`
    
    Results: 
    1. New Jira release is created, if it doesn't exist already - `1000.000.000`
    2. If task was found in commits and exists in Jira it will set fixed version for it
    3. If a task have been found in a commit messages, and it exists in Jira it is matched with the release.


## Contributing

### Linting
Project uses golint-ci as an aggrated linter.
You may find its configuration in [.golangci](https://github.com/g2a-com/kojiro/blob/master/.golangci.yaml) file
1. Lint: `make lint`
2. Lint with autofix: `make lint-fix`

### Versioning

We use [SemVer](http://semver.org/) for versioning to comply with golang standards.
For the versions available, see the [tags on this repository](https://github.com/g2a-com/kojiro/tags). 

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details


