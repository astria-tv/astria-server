# Hacking
This document is an overview of how to build and run Astria, as well as how to get your changes merged back upstream.

[[_TOC_]]

## Overview
The general idea is to check out a copy of the source code, compile it into a binary, and run it.

Once you're able to do that, you can start making changes to the code and test them. If you want, you can go and look at the [outstanding issues on GitHub](https://github.com/astria-tv/astria-server/issues/). There are tags for issues specifically targetted to new contributors [here](https://github.com/astria-tv/astria-server/issues?scope=all&utf8=%E2%9C%93&state=opened&label_name[]=Good%20first%20issue)

Once you've made a change in the code that fixes an issue (and tested it), the steps to get it accepted are [forking the repo on GitHub](https://docs.github.com/en/get-started/quickstart/fork-a-repo), pushing your change to a feature branch, and then submitting a [pull request](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request). 

### What's in this repo
The server handles the backend - scanning your media library, querying metadata from [The Movie DB](https://www.themoviedb.org/), building the database - as well as the frontend - serving the webapp via https, reading from the database and streaming videos from the backend to the frontend.

### What's not in this repo

The client (playback) software for any device that is not a modern web browser. ie. Android, iOS, Roku etc. There is an iOS client under development [here](https://github.com/astria-tv/astria-ios) but it's still in alpha state and a source-only release at this point.

### Supported platforms
 * Linux - these instructions were written for Arch but any distribution should work
 * Golang, v1.13+ recommended. You can install this through your distribution repos but the latest version is always available at https://golang.org
 * Git is required if you want to contribute changes back and optionally to obtain a copy of the code
 * Make
 * FFmpeg

# Building and running

## Getting the code
You can authenticate to GitHub using username/password or ssh. Setting up ssh keys is beyond the scope of this document but there are instructions [here](https://docs.github.com/en/authentication/connecting-to-github-with-ssh)

    git clone git@github.com:astria-tv/astria-server.git

## Install the toolchain

  * Install the Go toolchain, either from your distro repos or directly from the [Go project](https://golang.org/dl)
  * Install make
  * Install FFmpeg

## Build dependencies

There is a makefile that can handle various project tasks.

  * Run `make download-astria-react` to grab the latest build of the web frontend for Astria.

## Build astria

  * `make build-local` to build a binary for your local platform. The binary will be placed in `build/astria`.

## Running the server

By default, astria-server will open an existing database or create a new one and start listening for web connections on port 8080. For development you may want to override the defaults, for example to run against a copy of your primary astria db instead of the real one.

You can run the compiled binary as follows:
    `build/astria --config_dir ~/astria_dev_cfg/`

## Using Docker

Alternatively you can use the Docker development issue which supports hot-reloading of the codebase while you work on it. You can build the image using `docker build - < Dockerfile.dev -t astria-dev` from the astria-server repository. Once built you can run it using `docker run -p 8080:8080 -v $PWD:/go/src/github.com/astria-tv/astria-server -it astria-dev /go/bin/modd`

# Merging your changes upstream

So you've fixed a bug or added a new feature and now you want to merge your changes back to the main project so everyone can benefit.

## Before you start

### Ensure your code is formatted correctly
  * `go fmt -w <filename>` for each file you have modified

### Ensure tests are passing
You can test that the CI/CD pipeline will run locally as follows:
  *  `make vet`  runs a linter on the code
  *  `make test` to run the test suite ensure all tests still pass, including any new ones you've added

Adding tests is not required but is encouraged; making sure existing tests do not break is mandatory.

## Fork the repo and push your change
Once you've created an account on GitHub, [forking the repo](https://docs.github.com/en/get-started/quickstart/fork-a-repo) creates your own copy of the repo that you can push your changes to.

One possible git flow for pushing your changes is as follows:

  * `git remote set-url origin git@github.com:<your username>/astria-server.git` to point your local git client at your personal GitHub repo instead of the master Astria repo
  * `git checkout -b some-descriptive-name` to create a local feature branch
  * `git add <filenames>` to stage each file you've modified
  * `git commit -m 'some description of your work'` to commit your staged changes with a relevant message
  * `git push --set-upstream origin <your branch name>` to push your local branch to a branch on GitHub with the same name

## Submit a merge request
Once you have created a feature branch and pushed it to your fork of the repo, you can submit a [pull request](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request) to the Astria project. It never hurts to paste a link to your MR in the Discord channel but please be patient if nobody is able to look at it right away.

## Tips

After you've opened a merge request, you may see that the CI pipeline checks have passed. GitHub Actions will automatically run CI checks on your pull request.
