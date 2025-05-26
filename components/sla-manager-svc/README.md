
# SLA Manager service 

&copy; Atos Spain S.A. 2024

[![License: Apache v2](https://img.shields.io/badge/License-Apache%20v2-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)

----------------------------

## Description

The SLA Manager service is a Golang lightweight implementation of an SLA system, inspired by the WS-Agreement standard. Its features are:

* REST interface to manage the creation and deletion of agreements
* Agreements evaluation on background; any breach in the agreement terms generates an SLA violation
* Configurable monitoring: a monitoring has to be provided externally
* Configurable repository: a memory repository (for developing purposes) is provided, but more can be added.

### SLAs 

An agreement (SLA) is represented internally by a simple JSON structure (see more examples in resources/samples):

```json
{
    "id": "<SLA_ID>",
    "name": "<SLA_NAME>",
    "state": "started",
    "assessment": {},
    "creation": "2024-01-16T17:09:45Z",
    "expiration": "2026-01-16T17:09:45Z",
    "details": {
        "variables": [],
        "guarantees": [
            {
                "name": "<GUARANTEE_NAME>",
                "constraint": "<CONSTRAINT>"
            }
        ]
    }
}
```

----------------------------

## Usage guide

### Docker

#### Installation

Build the Docker image:

```bash
docker build -t sla-manager .

docker tag 76d68172d12e sla-manager:<version>
```

Run the container:

```bash
docker run -ti -p 8081:8080 sla-manager:latest
```
    
#### Configuration

Environment variables used by the SLA & QoS Manager:

  - **AGENT_ID** (e.g., "agente01")
  - **PROMETHEUS_ADDRESS** (e.g., "http://localhost:9090")
  - **MONITORING_ADAPTER** (e.g., "prometheus")
  - **NOTIFIER_ADAPTER** (e.g., "rest_endpoint", "rpc")
  - **NOTIFICATION_ENDPOINT** (e.g., "http://localhost:10090")
  - **CONTEXT_ZENOH_ENDPOINT** (e.g., "http://192.168.137.47:8000/dockerContextDefinitions/**")
  - **COMPOSE_PROJECT_NAME** (e.g., "sensor")
  
Run the container with env variables:

```bash
docker run -d -ti -e PROMETHEUS_ADDRESS='http://<IP>:9090' -e MONITORING_ADAPTER='prometheus' -p 8081:8080 sla-manager:latest

docker run -d -ti -e PROMETHEUS_ADDRESS='http://92.168.137.25:9090' -e MONITORING_ADAPTER='prometheus' -p 8081:8080 sla-manager:latest

docker logs a7b0ffd62e11 -f
```

#### Test applicatoin

##### Prometheus

```bash
docker run -d -ti -p 9090:9090 prom/prometheus:latest
```

##### Create SLA

```bash
curl -k -X POST -d @resources/service_definition_example_01.json 'http://<IP>:8081/api/v1/sla'

curl -k -X POST -d @resources/service_definition_example_01.json http://192.168.137.47:8081/api/v1/sla
```

Input example (service definition):

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
    "dockerRoleDefinitions": [
        {
            "id": "Processing",
            "imageId": "prhcatbsc/colmena-processing",
            "hardwareRequirements": [
                "CPU"
            ],
            "kpis": [{
                "query": "[go_memstats_frees_total] < 50000",
                "scope": "company_premises/building=."
            }]
        }
    ]
}
```


##### Delete SLA

```bash
curl -k -X DELETE -d @resources/sla_v1.json 'http://<IP>:8081/api/v1/sla/<SLA_ID>'

curl -k -X DELETE -d @resources/sla_v1.json http://192.168.137.25:8081/api/v1/sla/test_rest_service_01-eZYsQZh6bPWMGsMYUENsAZ
```

----------------------------

## LICENSES
SLA Manager component is licensed under [Apache License, version 2](LICENSE).