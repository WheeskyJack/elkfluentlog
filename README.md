# docker-compose example on instrumenting the go application for fluentd and elk stack
This Project is basic example on how to integrate go app logs, fluentd, amd elk stack using docker-compose.
Using this, user can play around with elastic search and can have tight control on the logs which needs to be injected.

## Pre-requisites

Ensure you install the latest version of make, docker and [docker-compose](https://docs.docker.com/compose/install/) on your Docker host machine. This has been tested with Docker for Mac.

# Quick Start

The code present is in main.go file. It generates predefined log lines in log file.
This log file is then in turn fed to the fluent contianer. This flient container then interacts with logstash and logstash in turn send these logs into elastic search. The init container in docker-compose creates datastream, index templates and ilm in elastic search.

To build and run the docker-compose :

1. (optional) add your code to main.go file
2. execute `make image` from top of repo
3. cd esdc
4. docker-compose up -d

Your docker-compose will be up and running. It takes about a minute to start the docker compose fully.
You can check elk stack using kibana at `http://localhost:5601` in your browser.

## Disclaimer:
Disclaimer: this setup is not in anyway secured. So use it only for inspiration.