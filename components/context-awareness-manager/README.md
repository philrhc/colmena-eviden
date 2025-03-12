# Context Awareness Manager

COLMENA Context Awareness Manager module handles the receipt, processing, and distribution of context for the service it want to be deployed. This module provides a REST API to receive contexts, sends notifications when a new context is received, and publishes these contexts to distributed network subscribers for further consumption.

## Additional Info

It integrates with Docker for dynamic container management and uses Zenoh for efficient context publication.

## Table of Contents

- [Project Structure](#project-structure)
- [Installation](#installation)
- [API Usage](#api-usage)
- [Swagger Documentation](#swagger-documentation)
- [Running Test](#running-test)
- [License](#license)

## Project Structure

The application mantain the following structure:

### Internal modules

#### Controller

The **controller** submodule is responsible for managing incoming HTTP requests, processing the received contexts, and communicating with other modules such as the database and monitor.

- **File**: [controllers.go](components/context-awareness-manager/internal/controllers/controllers.go)

It exposes endpoints for interacting with the system.

- **Endpoints**:

  - **POST /context**: Receives new context data in the form of a JSON payload and stores it in the database.
  - **GET /health**: Provides a simple health check to ensure that the service is running.

#### Database

The **database** submodule interacts with the underlying database, managing the creation of tables, inserting and retrieving context information, and maintaining persistent storage. It stores context data, which can be monitored and processed.

- **File**: [database.go](components/context-awareness-manager/internal/database/database.go)

### Monitor

The **monitor** submodule is responsible for the active monitoring of context values. It periodically checks for context changes, deploys Docker containers as needed, collects logs, and publishes these logs to a distributed network.

- **File**: [monitor.go](components/context-awareness-manager/internal/monitor/monitor.go)

### DockerClient

The **dockerclient** submodule provides the necessary functionality to interact with Docker containers. It allows the deployment of Docker containers, executes commands within them, collects logs, and cleans up the containers after use. This module abstracts Docker's API to simplify container management.

- **File**: [dockerclient.go](components/context-awareness-manager/internal/dockerclient/dockerclient.go)

### Common Packages

The **pkg** directory contains shared utility packages that are used across the project.

#### Logger

Handles logging across the application.

- **File**: [logger.go](components/context-awareness-manager/internal/pkg/logger/logger.go)

#### Response

Standardizes responses sent to the API clients.

- **File**: [response.go](components/context-awareness-manager/internal/pkg/response/response.go)

#### Server

Manages the HTTP server and routing.

- **File**: [server.go](components/context-awareness-manager/internal/pkg/server/server.go)

#### Models

Contains data models.

- **File**: [models.go](components/context-awareness-manager/internal/models/models.go)

Here is the Service Description model structure used in the application:

```go
package models

type ServiceDescription struct {
 ID                       ID                        `json:"id"`
 DockerContextDefinitions []DockerContextDefinition `json:"dockerContextDefinitions"`
 KPIs                     []string                  `json:"kpis"`
 DockerRoleDefinitions    []string                  `json:"dockerRoleDefinitions"`
}
```

## Installation

Ensure you have Go installed and set up on your machine. Initialize your Go module and get the necessary dependencies.

```sh
go mod tidy
```

### Deploy a Container

1. Clone the repository

    ```sh
    git clone <GITHUB_REPO>
    ```

2. Build the Docker image

    ```sh
    docker build -t <GITHUB_REPO>/context-awareness-manager -f components/context-awareness-manager/install/Dockerfile .
    ```

3. Run the Docker container

    ```sh
    docker compose -f deploy/compose/docker-compose.yaml up -d context-awareness-manager
    ```

4. Health check to the application

    Open your browser and go to <http://localhost:8080/health>

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
company_premises|jrubioc0/company_premises
sqlite> SELECT * FROM dockerRoleDefinitions;
Sensing|jrubioc0/colmena-sensing
Processing|jrubioc0/colmena-processing
sqlite> .exit
```

## API Usage

The application exposes a RESTful API for managing Docker contexts.

### POST /context

- **Endpoint**: `POST /context`
- **Description**: Adds a new Docker context to the system for monitoring.
- **Request Body**:

```json
{
    "id": {
        "value": "ExampleApplication"
    },
    "dockerContextDefinitions": [
        {
            "id": "company_premises",
            "imageId": "jrubioc0/company_premises"
        }
    ],
    "kpis": [],
    "dockerRoleDefinitions": []
}
```

- **Example Curl Command**

```sh
curl -X POST http://localhost:8080/context -H "Content-Type: application/json" -d '{
    "id": {
        "value": "ExampleApplication"
    },
    "dockerContextDefinitions": [
        {
            "id": "company_premises",
            "imageId": "jrubioc0/company_premises"
        }
    ],
    "kpis": [],
    "dockerRoleDefinitions": []
}'
```

## Swagger Documentation

To generate the Swagger documentation, run the following command in the root project folder:

```bash
swag init -g cmd/context-awareness-manager/main.go -o docs --parseDependency
go mod tidy
```

This command will create the Swagger documentation in the docs folder.

## Running Test

To run the unit tests for the application:

```sh
go fmt $(go list ./... | grep -v /vendor/)
go vet $(go list ./... | grep -v /vendor/)
go test $(go list ./... | grep -v /vendor/)
```

## Contribution

Tech:

- **Tech Stack**: Golang, Docker
- **CI/CD**: GitHub Actions (CICD Pipeline)
- **Databases**: MySQLite
- **Other Tools**: DockerClient, Zenoh, Go test (for testing)

Asset Owner:

- **Component Owner**: Maintained by the development team

A Gitflow methodology is implemented within this repo. Proceed as follows:

1. Open an Issue, label appropriately and assign to the next planned release.
2. Pull the latest develop and branch off of it with a branch or PR created from GitHub as draft.
3. Commit and push frequently.
4. When ready, set PR as ready, tag team and wait for approval.

## License

ATOS and Eviden Copyright applies. Apache License

```text
/*
 * Copyright 20XX-20XX, Atos Spain S.A.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *   * Redistributions of source code must retain the above copyright notice,
 *     this list of conditions and the following disclaimer.
 *   * Redistributions in binary form must reproduce the above copyright
 *     notice, this list of conditions and the following disclaimer in the
 *     documentation and/or other materials provided with the distribution.
 *   * Neither the name of the copyright holder nor the names of its
 *     contributors may be used to endorse or promote products derived from
 *     this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */
```
