# DEPLOYMENT TOOL (w/ PODMAN)

## Descripción

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
```

2. Install dependencies: `pip install -r requirements.txt`

2. Build the Docker image

```sh
docker build -t jrubioc0/deployment-tool -f components/deployment-tool/Dockerfile .
```

3. Run the Docker container

```sh
docker compose -f install/compose/docker-compose.yaml up -d context-awareness-manager
```

4. Access the application

Open your browser and go to http://localhost:8000/docs

## API Usage

### Sample Requests

- **Endpoint**: `POST /build-and-push`
- **Request Body**:

```json
{
  "base_directory": "/home/jrubioc/COLMENA/documentation/example_application/build",
  "repo_url": "local-registry:5000"
}
```
- **Example Curl Commmand**
```sh
curl -X POST http://localhost:8000/build_and_push -H "Content-Type: application/json" -d '{
  "base_directory": "/home/jrubioc/COLMENA/documentation/example_application/build",
  "repo_url": "local-registry:5000"
}'
```

## Swagger Documentation
To generate the Swagger documentation, annotate the controller methods and run the following command in the root project folder:

```bash
swag init -g app/rest_app.py -o docs
```
This command will create the Swagger documentation in the docs folder.

## Running Test

To run unit tests:

```sh
pytest tests/
```

## License
The Container Engine is released under the Apache 2.0 license.
Copyright © 2022-2024 Eviden. All rights reserved.
See the [LICENSE](LICENSE) file for more information.