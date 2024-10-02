# Docker Client Module

This module provides a function to deploy a Docker container by making an HTTP POST request to a microservice. The microservice deploys the container using the specified image and an optional command (`cmd`). The logs of the container are returned as a response.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Example](#example)

## Installation

Ensure you have Go installed and set up on your machine. Initialize your Go module and get the necessary dependencies.

```sh
go mod init containerengine
go mod tidy
```

## Usage

Import the module and use the deployContainer function to deploy a Docker container. You need to specify the Docker image and optionally a command (cmd) to run inside the container. If cmd is omitted, the default command of the image will be used.

## Example

```sh
curl -X POST http://localhost:8000/deploy -H "Content-Type: application/json" -d '{"image": "xaviercasasbsc/company_premises:latest"}'
```