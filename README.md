# Nexttrip
An application to provide timing on the arrival of public transit.

[![Go](https://github.com/cmarquis/nexttrip/actions/workflows/go.yml/badge.svg)](https://github.com/cmarquis/nexttrip/actions/workflows/go.yml)

## Development
The application is written in Go. See golang.org to setup a valid Go development environment.

The Cobra (https://github.com/spf13/cobra) module is used to handle CLI structure and arguments. As such main.go simply executes the cmd package.

### Packages

#### cmd
This is a Cobra structure application so the cmd package contains the cobra structs and functions to execute the application.

root.go is the entry point of the rootCmd and as of today only contains a single command.

#### providers
The providers package is where the communication with the transit providers lives. It uses dependency injection to allow for easy testing and mocking of transit provider API's.

### Testing
`go test ./...`

### Run
`go run ./...`

## Usage
The application can be compiled locally or a pre-built binary can be downloaded from GitHub releases.

### Usage:
  [ROUTE] [STOP] [DIRECTION] [flags]

### Flags:
  -h, --help   help for [ROUTE]