# govcms-cli
GovCMS's command line tool for local development.

# Requirements

- Git
- Docker

# Installation

## Homebrew

```shell
brew tap govcms-tests/govcms-cli
brew install govcms-cli
```

# Usage

Install local copies of the GovCMS distribution using:

```shell
govcms get
```

Build and launch GovCMS installations inside docker containers using:
```shell
govcms up <site_name>
```

For a more detailed view of the capabilities of the GovCMS CLI tool, check out the [docs]().
