# version [![godoc](https://img.shields.io/badge/go-documentation-blue.svg)](https://godoc.org/github.com/vaitekunas/version) [![Go Report Card](https://goreportcard.com/badge/github.com/vaitekunas/version)](https://goreportcard.com/report/github.com/vaitekunas/version) [![Build Status](https://travis-ci.org/vaitekunas/version.svg?branch=master)](https://travis-ci.org/vaitekunas/version) [![Coverage Status](https://coveralls.io/repos/github/vaitekunas/version/badge.svg?branch=master)](https://coveralls.io/github/vaitekunas/version?branch=master)

`version` is a small utility used to set and review semantic versions ([http://semver.org/](http://semver.org/))
of a git repository. It allows you to see already set versions and increase the most current version by a major/minor/patch/special
tick without the need to investigate which version should be set and if it does not violate the semver rules.

An example of a very simplified workflow (commiting code, increasing version and pushing it to origin, e.g github):

```shell
> git add .
> git commit -m "Fix something"
> version increase
> git push --tags origin master
```

`version` works directly with the git repository and does not require any additional files or configuration.

# Installing

`version` is written in [Go](https://golang.org) and requires the Go compiler to be installed:

``` shell
go get -u github.com/vaitekunas/version
go install github.com/vaitekunas/version
```

Assuming your `$PATH` environment variable includes `$GOPATH/bin`, you should be
able to run version in your shell:

```shell
> version --root=$GOPATH
```

# Using

`version` has only two methods:
* `version [--root] [--all]` - lists (`--all`) versions of all repositories (recursively) in the `--root` directory. If the root directory is not specified, then the working directory is used as root
* `version increase [{--major, --minor, --patch}] [--special=""] [--build=""]` - increases the version by a selected tick and sets it on the currently checked out/active commit.

Running version without any flags will list all the repositories (recursively, starting from the pwd)
and their highest version available:

```shell
> version

Highest versions per repository

╔════════════════════════════════════════════════════════╦══════════════════╦═════════╦═══════════════╗
║                       Repository                       ║       Date       ║ Commit  ║    Version    ║
╠════════════════════════════════════════════════════════╩══════════════════╩═════════╩═══════════════╣
║ /home/mindow/go/src/github.com/alecthomas/gometalinter │ 2017-03-22 00:27 │ bae2f12 │ v1.2.1        ║
╟────────────────────────────────────────────────────────┼──────────────────┼─────────┼───────────────╢
║ /home/mindow/go/src/github.com/derekparker/delve       │ 2017-05-06 01:09 │ f609169 │ v1.0.0-rc.1   ║
╟────────────────────────────────────────────────────────┼──────────────────┼─────────┼───────────────╢
║ /home/mindow/src/github.com/alecthomas/gometalinter    │ 2017-03-22 00:27 │ bae2f12 │ v1.2.1        ║
╟────────────────────────────────────────────────────────┼──────────────────┼─────────┼───────────────╢
║ /home/mindow/src/github.com/derekparker/delve          │ 2017-05-06 01:09 │ f609169 │ v1.0.0-rc.1   ║
╟────────────────────────────────────────────────────────┼──────────────────┼─────────┼───────────────╢
║ /home/mindow/src/github.com/fatih/color                │ 2017-05-23 15:53 │ 570b54c │ v1.5.0        ║
╟────────────────────────────────────────────────────────┼──────────────────┼─────────┼───────────────╢
║ /home/mindow/src/github.com/mattn/goveralls            │ 2016-09-12 16:13 │ 9d621f6 │ v0.0.1        ║
╟────────────────────────────────────────────────────────┼──────────────────┼─────────┼───────────────╢
║ /home/mindow/src/github.com/vaitekunas/version         │ 2017-09-07 16:02 │ 657b68d │ v0.1.0-rc1    ║
╟────────────────────────────────────────────────────────┼──────────────────┼─────────┼───────────────╢
║ /home/mindow/src/google.golang.org/grpc                │ 2017-03-14 23:44 │ cdee119 │ v1.2.0        ║
╟────────────────────────────────────────────────────────┼──────────────────┼─────────┼───────────────╢
║ /home/mindow/versailles                                │ 2017-09-07 16:09 │ 46a2962 │ v0.14.1       ║
╚════════════════════════════════════════════════════════╧══════════════════╧═════════╧═══════════════╝


───────────────────────────────────────────────────────────────────────────────────────
1. Version order is based on the semantic versioning specification (http://semver.org/)
2. Commits without version tags are not shown
```

Adding the `--all` flag lists all available versions/releases:

```shell
> cd ~/versailles
> version --all

All versions of '/home/mindow/versailles'
(ordered from the highest to the lowest)

╔═════════════════════════╦══════════════════╦═════════╦═════════════╗
║       Repository        ║       Date       ║ Commit  ║   Version   ║
╠═════════════════════════╩══════════════════╩═════════╩═════════════╣
║ /home/mindow/versailles │ 2017-09-07 16:09 │ 46a2962 │ v0.14.1     ║
╟─────────────────────────┼──────────────────┼─────────┼─────────────╢
║ /home/mindow/versailles │ 2017-09-07 16:05 │ 1788554 │ v0.14.0     ║
╟─────────────────────────┼──────────────────┼─────────┼─────────────╢
║ /home/mindow/versailles │ 2017-09-07 15:06 │ 9196116 │ v0.12.3     ║
╟─────────────────────────┼──────────────────┼─────────┼─────────────╢
║ /home/mindow/versailles │ 2017-09-07 12:15 │ 1f672c3 │ v0.12.2     ║
╟─────────────────────────┼──────────────────┼─────────┼─────────────╢
║ /home/mindow/versailles │ 2017-09-06 12:58 │ 2b11a34 │ v0.12.1     ║
╟─────────────────────────┼──────────────────┼─────────┼─────────────╢
║ /home/mindow/versailles │ 2017-09-06 15:54 │ abae609 │ v0.12.1-rc1 ║
╟─────────────────────────┼──────────────────┼─────────┼─────────────╢
║ /home/mindow/versailles │ 2017-09-06 10:40 │ 03331a3 │ v0.11.2     ║
╟─────────────────────────┼──────────────────┼─────────┼─────────────╢
║ /home/mindow/versailles │ 2017-09-06 10:36 │ 0a867e1 │ v0.11.1     ║
╚═════════════════════════╧══════════════════╧═════════╧═════════════╝


───────────────────────────────────────────────────────────────────────────────────────
1. Version order is based on the semantic versioning specification (http://semver.org/)
2. Commits without version tags are not shown
```

Running `version increase` without additional flags will propose a patch version update:

```shell
> version increase

Repository:
	 ◈  /home/mindow/versailles

Commit to be tagged as the new version:
	 ◈  Branch:  hotfix
	 ◈  Message: Fix bug#32
	 ◈  Hash:  6ba8e63f4433c7a1f6e7d0ffbc72c9c63ba06c34
	 ◈  Date:  2017-09-07 16:36:17
	 ◈  Author:  Mindaugas Vaitekunas

Version increment:
	 ◈  Current version: v0.14.1
	 ◈  Proposed version after increase: v0.14.2

Tag new version? [Y/n] (default: n):
```

if current version introduces breaking changes, then a major version update is
more appropriate:

```shell
> version increase --major

Repository:
	 ◈  /home/mindow/versailles

Commit to be tagged as the new version:
	 ◈  Branch:	hotfix
	 ◈  Message: Fix bug#32
	 ◈  Hash:	6ba8e63f4433c7a1f6e7d0ffbc72c9c63ba06c34
	 ◈  Date:	2017-09-07 16:36:17
	 ◈  Author:	Mindaugas Vaitekunas

Version increment:
	 ◈  Current version: v0.14.1
	 ◈  Proposed version after increase: v1.0.0

Tag new version? [Y/n] (default: n):
```

Ticks `--special` and `--build` can be used to set pre-release and build
versions respectively.

```shell
>  version increase --major --special="rc.1" --build="$(date +%s)"

Repository:
	 ◈  /home/mindow/versailles

Commit to be tagged as the new version:
	 ◈  Branch:	hotfix
	 ◈  Message:	Fix bug#32
	 ◈  Hash:	6ba8e63f4433c7a1f6e7d0ffbc72c9c63ba06c34
	 ◈  Date:	2017-09-07 16:36:17
	 ◈  Author:	Mindaugas Vaitekunas

Version increment:
	 ◈  Current version: v0.14.1
	 ◈  Proposed version after increase: v1.0.0-rc.1+1504795241

Tag new version? [Y/n] (default: n):
```

`version` *can* be combined with git hooks to increment versions automatically. Be
advised, however, that setting semantic versions will automatically create releases
on github, which is not necessarily what you want. Checking the branch before
incrementing the version would make more sense (e.g. increasing the version only when
commits are merged into the `release` branch).










# TODO

- [ ] Increase test coverage
