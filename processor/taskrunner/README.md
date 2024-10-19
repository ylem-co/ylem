# YLEM TASK RUNNER MICROSERVICE

![Static Badge](https://img.shields.io/badge/Go-1.23-black)
<a href="https://github.com/ylem-co/ylem?tab=Apache-2.0-1-ov-file">![Static Badge](https://img.shields.io/badge/license-Apache%202.0-black)</a>
<a href="https://ylem.co" target="_blank">![Static Badge](https://img.shields.io/badge/website-ylem.co-black)</a>
<a href="https://docs.ylem.co" target="_blank">![Static Badge](https://img.shields.io/badge/documentation-docs.ylem.co-black)</a>
<a href="https://join.slack.com/t/ylem-co/shared_invite/zt-2nawzl6h0-qqJ0j7Vx_AEHfnB45xJg2Q" target="_blank">![Static Badge](https://img.shields.io/badge/community-join%20Slack-black)</a>

Ylem Task Runners Microservice combines two functionalities:
1. Task runner that executes tasks and pipelines.
2. Load balancer that decides how to distribute tasks between Apache Kafka partitions.

More information can be found in our task-processing documentation: https://docs.ylem.co/open-source-edition/task-processing-architecture

# Usage

## Build an app

``` bash
$ go build
```

## Run database migrations

``` bash
$ ./ylem_taskrunner db migrations migrate
$ ./ylem_loadbalancer db migrations migrate
```

## Start application server

``` bash
$ ./ylem_taskrunner server serve
$ ./ylem_loadbalancer server serve
```

Task runner is now available inside the Ylem network on http://ylem_taskrunner:7335 or from the host machine on http://127.0.0.1:7335.

Load balancer is now available inside the Ylem network on http://ylem_loadbalancer:7334 or from the host machine on http://127.0.0.1:7334.

## Run tests

``` bash
$ go test ./tests/... -v -vet=off
```

## Run tests with test coverage

``` bash
$ go test ./tests/... -coverpkg=./... 
```

## Run tests with an advanced coverage report

``` bash
$ go test ./tests/... -coverpkg=./... -coverprofile cover.out
$ go tool cover -html cover.out -o cover.html
```

And then open cover.html

# Linter

## Install Golang linter on MacOS

``` bash
$ brew install golangci-lint
$ brew upgrade golangci-lint
```

## Check the code with it

``` bash
$ golangci-lint run
```

More information is in the official documentation: https://golangci-lint.run/
