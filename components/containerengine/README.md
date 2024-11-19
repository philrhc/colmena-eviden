# Container Engine

## Description

COLMENA Container Engine module belong to Context Awareness Component. This submodule provides a function to deploy a Docker container by making an HTTP POST request to a microservice. The microservice deploys the container using the specified image and an optional command (`cmd`). The logs of the container are returned as a response.

## Table of Contents

- [Installation](#installation)
- [API Usage](#api-usage)
- [Swagger Documentation](#swagger-documentation)
- [Running Test](#running-test)
- [License](#license)

## Project Structure

### Receiver

This submodule manages the reception of HTTP requests and processes the received contexts.

- **File**: [handler.go](components/containerengine/api/handlers/handler.go)

- **Endpoints**:
  - **POST /deploy**: This endpoint receives a new context image.
  - **GET /health**: This endpoint provides a simple health check.

### DockerClient

This module provides functionality to deploy Docker containers, execute commands within them, and retrieve their logs. It interacts with Docker's API to manage containers in a programmatic manner, ensuring easy deployment and log collection.

- **File**: [dockerclient.go](components/containerengine/internal/dockerclient/dockerclient.go)

### Models

The context managed within this module is defined in [models.go](components/containerengine/internal/models/models.go) and has the following structure:

```go
package models 

type DeployRequest struct {
    Image string   `json:"image"`
    Cmd   []string `json:"cmd,omitempty"`
}
```

## Installation

Ensure you have Go installed and set up on your machine. Initialize your Go module and get the necessary dependencies.

```sh
go mod init containerengine
go mod tidy
```

### Deploy a Container

1. Clone the repository

```sh
git clone https://github.com/eviden-colmena/colmena-eviden.git
```

2. Build the Docker image

```sh
docker build -t jrubioc0/containerengine -f components/containerengine/build/Dockerfile .
```

3. Run the Docker container

```sh
docker compose -f install/compose/docker-compose.yaml up -d containerengine
```

4. Access the application

Open your browser and go to http://localhost:8080/health

## API usage

### Sample Requests

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
- **Example Curl Commmand**
```sh
curl -X POST http://localhost:8000/deploy \
    -H "Content-Type: application/json" \ 
    -d '{"image": "xaviercasasbsc/company_premises:latest"}'
```

## Swagger Documentation
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