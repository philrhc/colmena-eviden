# Context Awareness Manager

## Descripción

El módulo Gestor de Conciencia de Contexto es una aplicación desarrollada en Go que se encarga de recibir, procesar y distribuir contextos a través de diferentes submódulos. Este módulo proporciona una API REST para recibir contextos, envía notificaciones cuando un nuevo contexto es recibido y publica estos contextos a los suscriptores de la red Zenoh para que puedan consumirlos.

## Submódulos

### Context

Submódulo que maneja la lógica del contexto.

- **Archivo**: [main.go](context-awareness-manager-svc/main.go)

- **Endpoints**:
  - **PUT /context**: Este endpoint recibe un contexto nuevo.
  - **GET /health**: Este endpoint proporciona una verificación de estado simple.
  - **GET /**: Este endpoint proporciona un mensaje de bienvenida.

### Receiver

Submódulo que gestiona la recepción de la petición HTTP y procesa el contexto recibido.

- **Archivo**: [receiver/handler.go](context-awareness-manager-svc/src/receiver/handler.go)

### Monitor

Submódulo que interactua con los microservicios encargados de obtener el contexto y la red Zenoh mediante publicaciones HTTP POST cuando se recibe un nuevo contexto.

- **Archivo**: [monitor/monitor.go](context-awareness-manager-svc/src/monitor/monitor.go)

## Estructura del Contexto

El contexto que se maneja en este módulo está definido en el archivo [models.go](context-awareness-manager-svc/src/models.go) y tiene la siguiente estructura:

```go
package context

type Context struct {
	ID                       ID                        `json:"id"`
	DockerContextDefinitions []DockerContextDefinition `json:"dockerContextDefinitions"`
	KPIs                     []KPI                     `json:"kpis"`
	DockerRoleDefinitions    []DockerRoleDefinition    `json:"dockerRoleDefinitions"`
}
```

## Pasos para construir y ejecutar la imagen Docker

1. Construir la imagen Docker

```sh
docker build -t context-awareness-manager -f context-awareness-manager-svc/Dockerfile .
```

2. Ejecutar el contenedor Docker

```sh
docker run -d \
  --name context-awareness \
  --restart always \
  -p 8080:8080 \
  -e DOCKERENGINE_URL="http://containerengine:9000/deploy" \
  registry.atosresearch.eu:18512/context-awareness-manager:develop
```

3. Acceder a la aplicación

Abre tu navegador y visita http://localhost:8080.

4. Ver los datos en la base de datos SQLite

```sh
$ sqlite3 ./context-awareness-manager-svc/context_awareness_manager.db
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

### Ejemplo de Solicitudes

# Enviar un contexto:
Usar script [request_service_decription](documentation/resources/request.sh)

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
    "dockerRoleDefinitions": [
        {
            "id": "Sensing",
            "imageId": "xaviercasasbsc/colmena-sensing",
            "hardwareRequirements": [
                "CAMERA"
            ],
            "kpis": []
        },
        {
            "id": "Processing",
            "imageId": "xaviercasasbsc/colmena-processing",
            "hardwareRequirements": [
                "CPU"
            ],
            "kpis": [
                {
                    "value": "buffer_queue_size[100000000s] < 10"
                }
            ]
        }
    ]
}'
```
