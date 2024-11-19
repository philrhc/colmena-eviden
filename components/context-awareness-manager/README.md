# Context Awareness Manager

## Descripción

COLMENA Context Awareness Manager module handles the receipt, processing, and distribution of context for the service it want to be deployed. This module provides a REST API to receive contexts, sends notifications when a new context is received, and publishes these contexts to distributed network subscribers for further consumption.

## Table of Contents

- [Project Structure](#project-structure)
- [Installation](#installation)
- [API Usage](#api-usage)
- [Swagger Documentation](#swagger-documentation)
- [Running Test](#running-test)
- [License](#license)

## Project Structure

### Receiver

This submodule manages the reception of HTTP requests and processes the received contexts.

- **File**: [handler.go](components/context-awareness-manager/api/handlers/handler.go)

- **Endpoints**:
  - **POST /context**: This endpoint receives a new context.
  - **GET /health**: This endpoint provides a simple health check.

### Monitor

This submodule monitoring context values and interacts with others microservices to receive this information and notify to al distributed modules of the system.responsible for obtaining context value and with the distributed network, using HTTP PUT requests to publish new contexts.

- **File**: [monitor.go](components/context-awareness-manager/internal/monitor/monitor.go)

- **Microservices involved**:
  - **ContainerEngine**: Context Awareness Manager submodule responsible that processes context images and return their logs.
  - **ZenohRouter**: Agent component acting as an endpoint within the Distributed Network.

### Models

The context managed within this module is defined in [models.go](components/context-awareness-manager/internal/models/models.go) and has the following structure:

```go
package context

type Context struct {
	ID                       ID                        `json:"id"`
	DockerContextDefinitions []DockerContextDefinition `json:"dockerContextDefinitions"`
	KPIs                     []KPI                     `json:"kpis"`
	DockerRoleDefinitions    []DockerRoleDefinition    `json:"dockerRoleDefinitions"`
}
```

## Installation

Ensure you have Go installed and set up on your machine. Initialize your Go module and get the necessary dependencies.

```sh
go mod init context-awareness-manager
go mod tidy
```

### Deploy a Container

1. Clone the repository

```sh
git clone https://github.com/eviden-colmena/colmena-eviden.git
```

2. Build the Docker image

```sh
docker build -t jrubioc0/context-awareness-manager -f components/context-awareness-manager/build/Dockerfile .
```

3. Run the Docker container

```sh
docker compose -f install/compose/docker-compose.yaml up -d context-awareness-manager
```

4. Access the application

Open your browser and go to http://localhost:8080/health

5. View data in the SQLite database:

```sh
$ sqlite3 ./components/context-awareness-manager/context_awareness_manager.db
SQLite version 3.34.1 2021-01-20 14:10:07
Enter ".help" for usage hints.
sqlite> .tables
dockerContextDefinitions
sqlite> .schema dockerContextDefinitions
CREATE TABLE IF NOT EXISTS dockerContextDefinitions (
    id TEXT PRIMARY KEY,
    imageId TEXT NOT NULL
);
sqlite> SELECT * FROM dockerContextDefinitions;
company_premises|xaviercasasbsc/company_premises
sqlite> SELECT * FROM dockerRoleDefinitions;
Sensing|xaviercasasbsc/colmena-sensing
Processing|xaviercasasbsc/colmena-processing
sqlite> .exit
```

## API Usage

### Sample Requests

- **Endpoint**: `POST /context`
- **Request Body**:

```json
{
    "id": {
        "value": "ExampleApplication"
    },
    "dockerContextDefinitions": [
        {
            "id": "company_premises",
            "imageId": "xaviercasasbsc/company_premises"
        }
    ],
    "kpis": [],
    "dockerRoleDefinitions": []
}
```
- **Example Curl Commmand**
```sh
curl -X POST http://localhost:8080/context -H "Content-Type: application/json" -d '{
    "id": {
        "value": "ExampleApplication"
    },
    "dockerContextDefinitions": [
        {
            "id": "company_premises",
            "imageId": "xaviercasasbsc/company_premises"
        }
    ],
    "kpis": [],
    "dockerRoleDefinitions": []
}'
```

## Swagger Documentation
To generate the Swagger documentation, annotate the controller methods and run the following command in the root project folder:

```bash
swag init -g cmd/context-awareness-manager/main.go -o docs
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