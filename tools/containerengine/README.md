# Docker Client Module

This module provides a function to deploy a Docker container by making an HTTP POST request to a microservice. The microservice deploys the container using the specified image and an optional command (`cmd`). The logs of the container are returned as a response.

## Table of Contents

- [Installation](#installation)
- [API Usage](#api-usage)
- [Deploy a Container](#deploy-a-container)
- [Generating Swagger Documentation](#generating-swagger-documentation)
- [Running Test](#running-test)
- [License](#license)

## Installation

Ensure you have Go installed and set up on your machine. Initialize your Go module and get the necessary dependencies.

```sh
go mod init containerengine
go mod tidy
```

## API usage

Import the module and use the deployContainer function to deploy a Docker container. You need to specify the Docker image and optionally a command (cmd) to run inside the container. If cmd is omitted, the default command of the image will be used.

#### Deploy a Container

- **Endpoint**: `POST /deploy`
- **Request Body**:

```json
{
  "image": "xaviercasasbsc/company_premises:latest"
}
```
- **Response Body**:
```json
{
  "classification": "reception"
}
```
** Example Curl Commmand**
```sh
curl -X POST http://localhost:8000/deploy \
    -H "Content-Type: application/json" \ 
    -d '{"image": "xaviercasasbsc/company_premises:latest"}'
```

## Generating Swagger Documentation
To generate the Swagger documentation, annotate the controller methods and run the following command in the root project folder:

```bash
swag init -g cmd/containerengine/main.go -o docs
```
This command will create the Swagger documentation in the docs folder.

## Running Test

To run unit tests:

```sh
go test ./..
```

## License
The Container Engine is released under the Apache 2.0 license.
Copyright © 2022-2024 Eviden. All rights reserved.
See the [LICENSE](LICENSE) file for more information.