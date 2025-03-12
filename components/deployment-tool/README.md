# DEPLOYMENT TOOL

COLMENA Deployment Tool provides an API to build and push Docker images for distributed modules of a service provide in a specified base path. Also publish the service description to distributed network subscribers for further consumption.

## Table of Contents

- [Project Structure](#project-structure)
- [Installation](#installation)
- [API Usage](#api-usage)
- [Swagger Documentation](#swagger-documentation)
- [Running Test](#running-test)
- [License](#license)

## Project Structure

- **Build and push Docker Image:**
  - **Propósito:** Construye imagenes Docker en base a unos builds artifacts.
  - **Funcionamiento:** Escucha peticiones (POST) del Panel de Control donde recibe un path directory donde se encuentran los builds artifacts y el service description. En base a esto construye las imagenes Docker y las pushea a Docker Hub.

- **Publish Service Description:**
  - **Propósito:** Envía una publicacion con el Service Description obtenido a la red Zenoh.
  - **Funcionamiento:** Escucha peticiones (POST) del Panel de Control donde recibe un path directory donde se encuentran los builds artifacts y el service description. En base a esto publica en la red Zenoh, con key_expression "service_description" el contenido del Service Description.

## Installation

### Deploy a Container

1. Clone this repository

    ```sh
    git clone https://github.com/eviden-colmena/colmena-eviden.git
    cd components/deployment-tool
    ```

2. Install dependencies

    ```bash
    python -m venv venv
    source venv/bin/activate
    pip install -r requirements.txt
    # exit venv with `deactivate`
    ```

3. Build the Docker image

    ```sh
    docker build -t <GITHUB_REPO>/deployment-tool -f components/deployment-tool/install/Dockerfile .
    ```

4. Run the Docker container

    ```sh
    docker compose -f install/compose/docker-compose.yaml up -d deployment-tool
    ```

5. Access the application

Open your browser and go to <http://localhost:8000/docs>

## API Usage

### Sample Requests

- **Endpoint**: `GET /deployservice`

- **Example Curl Command**

```sh
curl http://localhost:8000/deployservice
```

- **Response body**

```sh
{
  "message": "Images processed and service description published."
}
```

- **Validate image creation (local-registry)**

```sh
  curl -u <username>:<password> <http://localhost:5000/v2/_catalog>
```

## Swagger Documentation

Available at run time in <http://127.0.0.1:8000/docs>

## Running Test

To run unit tests:

```sh
pytest tests --cov=app --cov-report html --cov-report term-missing >> report.txt
pylint app --exit-zero --reports y >> qa_report.txt
```

## Contribution

Tech:

- **Tech Stack**: FastAPI, Python, Docker, Kubernetes, Helm
- **CI/CD**: GitHub Actions (CICD Pipeline)
- **Databases**: Not applicable for this component
- **Other Tools**: Docker, Uvicorn, pytest (for testing)

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
