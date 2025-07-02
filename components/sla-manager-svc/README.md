
# SLA Manager service 

&copy; Atos Spain S.A. 2024

[![License: Apache v2](https://img.shields.io/badge/License-Apache%20v2-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)

----------------------------

- [SLA Manager service](#sla-manager-service)
  - [1. Description](#1-description)
    - [SLAs](#slas)
  - [2. Usage guide](#2-usage-guide)
    - [2.1 Installation](#21-installation)
    - [2.2 Configuration](#22-configuration)
    - [2.3 Test application](#23-test-application)
        - [Required](#required)
        - [Create SLA](#create-sla)
          - [POST api/v1/sla](#post-apiv1sla)
          - [GET api/v1/slas](#get-apiv1slas)
          - [GET api/v1/slas/:id](#get-apiv1slasid)
          - [GET api/v1/sla/:id](#get-apiv1slaid)
        - [Delete SLA](#delete-sla)
  - [3. SLAs with scope](#3-slas-with-scope)
    - [Service descriptor](#service-descriptor)
    - [PAUSED SLA](#paused-sla)
    - [Send CONTEXT and METRICS to PROMETHEUS](#send-context-and-metrics-to-prometheus)
    - [STARTED SLA](#started-sla)
      - [ASSESSMENT OK](#assessment-ok)
      - [VIOLATION](#violation)
  - [4. KPI queries](#4-kpi-queries)
  - [5. Notifications and violations](#5-notifications-and-violations)
      - [NOTIFICATION](#notification)
      - [VIOLATIONs](#violations)
  - [LICENSES](#licenses)

----------------------------

## 1. Description

The SLA Manager service is a Golang lightweight implementation of an SLA system, inspired by the WS-Agreement standard. Its features are:

* REST interface to manage the creation and deletion of agreements
* Agreements evaluation on background; any breach in the agreement terms generates an SLA violation
* Configurable monitoring: a monitoring has to be provided externally. In the case of COLMENA a Prometheus instance will be used.
* Configurable repository: a memory repository (for developing purposes) is provided, but more can be added.

The SLA Manager provides the following methods:

- **POST api/v1/sla** creates a SLA
- **GET api/v1/sla/:id** gets the information about a specific SLA
- **DELETE api/v1/sla/:id** deletes a SLA
- **GET api/v1/slas** gets the information about all the SLAs
- **GET api/v1/slas/:id** gets the information about all the SLAs of a specific service
- **GET api/v1/kpis** gets the information about all the SLAs (KPI format)
- **GET api/v1/kpis/:id** gets the information about all the SLAs of a specific service (KPI format)
- **GET api/v1/kpi/:id** gets the information about a specific SLA (KPI format)

### SLAs 

An agreement (SLA) is represented internally by a simple JSON structure (see more examples in resources/samples):

```json
{
    "id": "<SLA_IDENTIFIER>",
    "name": "<SERVICE_NAME>",
    "state": "started",
    "assessment": {
      "total_executions": 0,
      "total_violations": 0,
      "level": "<LEVEL>",
      "violated": false,
      "first_execution": "2025-05-07T10:15:41.741141409Z",
      "last_execution": "2025-05-12T13:53:41.74250365Z",
      "guarantees": {
        "Processing01": {
          "first_execution": "2025-05-07T10:15:41.741141409Z",
          "last_execution": "2025-05-12T13:53:41.74250365Z",
          "last_values": {
            "": {
              "key": "",
              "action": "",
              "namespace": "",
              "value": 0,
              "datetime": "2025-05-12T13:53:41.743Z"
            }
          },
          "last_violation": {
            "id": "",
            "agreement_id": "<SLA_IDENTIFIER>",
            "guarantee": "",
            "datetime": "2025-05-12T13:53:41.743Z",
            "constraint": "",
            "values": [
              {
                "key": "",
                "action": "",
                "namespace": "",
                "value": 0,
                "datetime": "2025-05-12T13:53:41.743Z"
              }
            ],
            "appID": ""
          }
        }
      }
    },
    "creation": "2025-05-07T10:15:41.314266531Z",
    "expiration": "2026-05-07T10:15:41.31426687Z",
    "details": {
      "guarantees": [
        {
          "name": "<GUARANTEE_NAME>",
          "constraint": "<CONSTRAINT_QUERY>",
          "query": "<CONSTRAINT_QUERY_TEMPLATE>",
          "scope": "",
          "scopeTemplate": ""
        }
      ]
    }
  }
```

----------------------------

## 2. Usage guide

### 2.1 Installation

You can download the code and run the applications after creating the Docker image.

Build the Docker image:

```bash
docker build -t sla-manager .

docker tag 76d68172d12e sla-manager:<version>
```

Run the container:

```bash
docker run -ti -p 8081:8080 sla-manager:latest
```
    
### 2.2 Configuration

The following environment variables are used by the SLA & QoS Manager:

  - Prometheus / Local Metric collector:
    - **PROMETHEUS_ADDRESS** (e.g., "http://prometheus:9090")
    - **MONITORING_ADAPTER** (e.g., "prometheus")
  - Notifications / Violations:
    - **NOTIFIER_ADAPTER** (e.g., "rest_endpoint")
    - **NOTIFICATION_ENDPOINT** (e.g., "http://localhost:10090")
  - Zenoh:
    - **CONTEXT_ZENOH_ENDPOINT** (e.g., "http://zenoh-router:8000")
    - **CONTEXT_ZENOH_CONTEXTS** (e.g., "colmena/contexts")
  - Agent Identifier: **COMPOSE_PROJECT_NAME** or **AGENT_ID** (e.g., "sensor", "ColmenaAgent1")
  
### 2.3 Test application

##### Required

- Prometheus is running and listening in port 9090 (set value in environment variable)
- Zenoh-Prometheus connector (metrics-etl) is running and listening in 8999

##### Create SLA
###### POST api/v1/sla

SLAs are created with the information provided in the **service descriptors**:

```bash
curl -k -X POST -d @resources/service_definition_example_01.json 'http://<IP>:8081/api/v1/sla'
```

Input example (service descriptor):

```json
{
    "id": {
        "value": "ExampleApplication_01"
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
            "id": "Processing01",
            "imageId": "prhcatbsc/colmena-processing",
            "hardwareRequirements": [
                "CPU"
            ],
            "kpis": [{
                "query": "[go_memstats_frees_total] < 50000",
                "scope": ""
            }]
        },
        {
            "id": "Processing02",
            "imageId": "prhcatbsc/colmena-processing",
            "hardwareRequirements": [
                "CPU"
            ],
            "kpis": [{
                "query": "[go_memstats_frees_total] < 50000",
                "scope": ""
            }]
        },
        {
            "id": "Processing03",
            "imageId": "prhcatbsc/colmena-processing",
            "hardwareRequirements": [
                "CPU"
            ],
            "kpis": [{
                "query": "[go_memstats_frees_total] < 50000",
                "scope": ""
            }]
        }
    ]
}
```

The previous service descriptor creates the following three SLAs:

###### GET api/v1/slas

e.g. http://sla-manager:8081/api/v1/slas

```json
{
  "Message": "Objects found",
  "Method": "GetSLAs",
  "Resp": "ok",
  "Response": [
    {
      "id": "ExampleApplication_01-Enw6R5Pni7eanXVHtEM8sR",
      "name": "ExampleApplication_01",
      "state": "started",
      "assessment": {
        "total_executions": 14834,
        "total_violations": 14834,
        "x": 2,
        "x_assessment_broken_count": 14834,
        "y": 2,
        "z": 5,
        "level": "Critical",
        "violated": true,
        "first_execution": "2025-05-07T10:15:41.741141409Z",
        "last_execution": "2025-05-12T13:52:11.741001441Z",
        "guarantees": {
          "Processing01": {
            "first_execution": "2025-05-07T10:15:41.741141409Z",
            "last_execution": "2025-05-12T13:52:11.741001441Z",
            "last_values": {
              "go_memstats_frees_total": {
                "key": "go_memstats_frees_total",
                "action": "",
                "namespace": "",
                "value": 239641613,
                "datetime": "2025-05-12T13:52:11.741Z"
              }
            },
            "last_violation": {
              "id": "",
              "agreement_id": "ExampleApplication_01-Enw6R5Pni7eanXVHtEM8sR",
              "guarantee": "Processing01",
              "datetime": "2025-05-12T13:52:11.741Z",
              "constraint": "[go_memstats_frees_total] \u003C 50000",
              "values": [
                {
                  "key": "go_memstats_frees_total",
                  "action": "",
                  "namespace": "",
                  "value": 239641613,
                  "datetime": "2025-05-12T13:52:11.741Z"
                }
              ],
              "appID": "ExampleApplication_01-Enw6R5Pni7eanXVHtEM8sR"
            }
          }
        }
      },
      "creation": "2025-05-07T10:15:41.314266531Z",
      "expiration": "2026-05-07T10:15:41.31426687Z",
      "details": {
        "guarantees": [
          {
            "name": "Processing01",
            "constraint": "[go_memstats_frees_total] \u003C 50000",
            "query": "[go_memstats_frees_total#LABELS#] \u003C 50000",
            "scope": "",
            "scopeTemplate": ""
          }
        ]
      }
    },
    {
      "id": "ExampleApplication_01-SA4H2HxucWoE9RsWg4AYGY",
      "name": "ExampleApplication_01",
      "state": "started",
      "assessment": {
        "total_executions": 14834,
        "total_violations": 14834,
        "x": 2,
        "x_assessment_broken_count": 14834,
        "y": 2,
        "z": 5,
        "level": "Critical",
        "violated": true,
        "first_execution": "2025-05-07T10:15:41.741141409Z",
        "last_execution": "2025-05-12T13:52:11.741001441Z",
        "guarantees": {
          "Processing02": {
            "first_execution": "2025-05-07T10:15:41.741141409Z",
            "last_execution": "2025-05-12T13:52:11.741001441Z",
            "last_values": {
              "go_memstats_frees_total": {
                "key": "go_memstats_frees_total",
                "action": "",
                "namespace": "",
                "value": 239641613,
                "datetime": "2025-05-12T13:52:11.742Z"
              }
            },
            "last_violation": {
              "id": "",
              "agreement_id": "ExampleApplication_01-SA4H2HxucWoE9RsWg4AYGY",
              "guarantee": "Processing02",
              "datetime": "2025-05-12T13:52:11.742Z",
              "constraint": "[go_memstats_frees_total] \u003C 50000",
              "values": [
                {
                  "key": "go_memstats_frees_total",
                  "action": "",
                  "namespace": "",
                  "value": 239641613,
                  "datetime": "2025-05-12T13:52:11.742Z"
                }
              ],
              "appID": "ExampleApplication_01-SA4H2HxucWoE9RsWg4AYGY"
            }
          }
        }
      },
      "creation": "2025-05-07T10:15:41.314304476Z",
      "expiration": "2026-05-07T10:15:41.314304768Z",
      "details": {
        "guarantees": [
          {
            "name": "Processing02",
            "constraint": "[go_memstats_frees_total] \u003C 50000",
            "query": "[go_memstats_frees_total#LABELS#] \u003C 50000",
            "scope": "",
            "scopeTemplate": ""
          }
        ]
      }
    },
    {
      "id": "ExampleApplication_01-kJRBWUF8wDMr4qx3FQ4R3q",
      "name": "ExampleApplication_01",
      "state": "started",
      "assessment": {
        "total_executions": 14834,
        "total_violations": 14834,
        "x": 2,
        "x_assessment_broken_count": 14834,
        "y": 2,
        "z": 5,
        "level": "Critical",
        "violated": true,
        "first_execution": "2025-05-07T10:15:41.741141409Z",
        "last_execution": "2025-05-12T13:52:11.741001441Z",
        "guarantees": {
          "Processing03": {
            "first_execution": "2025-05-07T10:15:41.741141409Z",
            "last_execution": "2025-05-12T13:52:11.741001441Z",
            "last_values": {
              "go_memstats_frees_total": {
                "key": "go_memstats_frees_total",
                "action": "",
                "namespace": "",
                "value": 239641613,
                "datetime": "2025-05-12T13:52:11.744Z"
              }
            },
            "last_violation": {
              "id": "",
              "agreement_id": "ExampleApplication_01-kJRBWUF8wDMr4qx3FQ4R3q",
              "guarantee": "Processing03",
              "datetime": "2025-05-12T13:52:11.744Z",
              "constraint": "[go_memstats_frees_total] \u003C 50000",
              "values": [
                {
                  "key": "go_memstats_frees_total",
                  "action": "",
                  "namespace": "",
                  "value": 239641613,
                  "datetime": "2025-05-12T13:52:11.744Z"
                }
              ],
              "appID": "ExampleApplication_01-kJRBWUF8wDMr4qx3FQ4R3q"
            }
          }
        }
      },
      "creation": "2025-05-07T10:15:41.314540964Z",
      "expiration": "2026-05-07T10:15:41.314541344Z",
      "details": {
        "guarantees": [
          {
            "name": "Processing03",
            "constraint": "[go_memstats_frees_total] \u003C 50000",
            "query": "[go_memstats_frees_total#LABELS#] \u003C 50000",
            "scope": "",
            "scopeTemplate": ""
          }
        ]
      }
    }
  ]
}
```

###### GET api/v1/slas/:id

e.g. http://sla-manager:8081/api/v1/slas/ExampleApplication_01

```json
{
  "Message": "Object found",
  "Method": "GetSLAsByServiceId",
  "Resp": "ok",
  "Response": [
    {
      "id": "ExampleApplication_01-kJRBWUF8wDMr4qx3FQ4R3q",
      "name": "ExampleApplication_01",
      "state": "started",
      "assessment": {
        "total_executions": 14837,
        "total_violations": 14837,
        "x": 2,
        "x_assessment_broken_count": 14837,
        "y": 2,
        "z": 5,
        "level": "Critical",
        "violated": true,
        "first_execution": "2025-05-07T10:15:41.741141409Z",
        "last_execution": "2025-05-12T13:53:41.74250365Z",
        "guarantees": {
          "Processing03": {
            "first_execution": "2025-05-07T10:15:41.741141409Z",
            "last_execution": "2025-05-12T13:53:41.74250365Z",
            "last_values": {
              "go_memstats_frees_total": {
                "key": "go_memstats_frees_total",
                "action": "",
                "namespace": "",
                "value": 239692000,
                "datetime": "2025-05-12T13:53:41.746Z"
              }
            },
            "last_violation": {
              "id": "",
              "agreement_id": "ExampleApplication_01-kJRBWUF8wDMr4qx3FQ4R3q",
              "guarantee": "Processing03",
              "datetime": "2025-05-12T13:53:41.746Z",
              "constraint": "[go_memstats_frees_total] \u003C 50000",
              "values": [
                {
                  "key": "go_memstats_frees_total",
                  "action": "",
                  "namespace": "",
                  "value": 239692000,
                  "datetime": "2025-05-12T13:53:41.746Z"
                }
              ],
              "appID": "ExampleApplication_01-kJRBWUF8wDMr4qx3FQ4R3q"
            }
          }
        }
      },
      "creation": "2025-05-07T10:15:41.314540964Z",
      "expiration": "2026-05-07T10:15:41.314541344Z",
      "details": {
        "guarantees": [
          {
            "name": "Processing03",
            "constraint": "[go_memstats_frees_total] \u003C 50000",
            "query": "[go_memstats_frees_total#LABELS#] \u003C 50000",
            "scope": "",
            "scopeTemplate": ""
          }
        ]
      }
    },
    {
      "id": "ExampleApplication_01-Enw6R5Pni7eanXVHtEM8sR",
      "name": "ExampleApplication_01",
      "state": "started",
      "assessment": {
        "total_executions": 14837,
        "total_violations": 14837,
        "x": 2,
        "x_assessment_broken_count": 14837,
        "y": 2,
        "z": 5,
        "level": "Critical",
        "violated": true,
        "first_execution": "2025-05-07T10:15:41.741141409Z",
        "last_execution": "2025-05-12T13:53:41.74250365Z",
        "guarantees": {
          "Processing01": {
            "first_execution": "2025-05-07T10:15:41.741141409Z",
            "last_execution": "2025-05-12T13:53:41.74250365Z",
            "last_values": {
              "go_memstats_frees_total": {
                "key": "go_memstats_frees_total",
                "action": "",
                "namespace": "",
                "value": 239692000,
                "datetime": "2025-05-12T13:53:41.743Z"
              }
            },
            "last_violation": {
              "id": "",
              "agreement_id": "ExampleApplication_01-Enw6R5Pni7eanXVHtEM8sR",
              "guarantee": "Processing01",
              "datetime": "2025-05-12T13:53:41.743Z",
              "constraint": "[go_memstats_frees_total] \u003C 50000",
              "values": [
                {
                  "key": "go_memstats_frees_total",
                  "action": "",
                  "namespace": "",
                  "value": 239692000,
                  "datetime": "2025-05-12T13:53:41.743Z"
                }
              ],
              "appID": "ExampleApplication_01-Enw6R5Pni7eanXVHtEM8sR"
            }
          }
        }
      },
      "creation": "2025-05-07T10:15:41.314266531Z",
      "expiration": "2026-05-07T10:15:41.31426687Z",
      "details": {
        "guarantees": [
          {
            "name": "Processing01",
            "constraint": "[go_memstats_frees_total] \u003C 50000",
            "query": "[go_memstats_frees_total#LABELS#] \u003C 50000",
            "scope": "",
            "scopeTemplate": ""
          }
        ]
      }
    },
    {
      "id": "ExampleApplication_01-SA4H2HxucWoE9RsWg4AYGY",
      "name": "ExampleApplication_01",
      "state": "started",
      "assessment": {
        "total_executions": 14837,
        "total_violations": 14837,
        "x": 2,
        "x_assessment_broken_count": 14837,
        "y": 2,
        "z": 5,
        "level": "Critical",
        "violated": true,
        "first_execution": "2025-05-07T10:15:41.741141409Z",
        "last_execution": "2025-05-12T13:53:41.74250365Z",
        "guarantees": {
          "Processing02": {
            "first_execution": "2025-05-07T10:15:41.741141409Z",
            "last_execution": "2025-05-12T13:53:41.74250365Z",
            "last_values": {
              "go_memstats_frees_total": {
                "key": "go_memstats_frees_total",
                "action": "",
                "namespace": "",
                "value": 239692000,
                "datetime": "2025-05-12T13:53:41.744Z"
              }
            },
            "last_violation": {
              "id": "",
              "agreement_id": "ExampleApplication_01-SA4H2HxucWoE9RsWg4AYGY",
              "guarantee": "Processing02",
              "datetime": "2025-05-12T13:53:41.744Z",
              "constraint": "[go_memstats_frees_total] \u003C 50000",
              "values": [
                {
                  "key": "go_memstats_frees_total",
                  "action": "",
                  "namespace": "",
                  "value": 239692000,
                  "datetime": "2025-05-12T13:53:41.744Z"
                }
              ],
              "appID": "ExampleApplication_01-SA4H2HxucWoE9RsWg4AYGY"
            }
          }
        }
      },
      "creation": "2025-05-07T10:15:41.314304476Z",
      "expiration": "2026-05-07T10:15:41.314304768Z",
      "details": {
        "guarantees": [
          {
            "name": "Processing02",
            "constraint": "[go_memstats_frees_total] \u003C 50000",
            "query": "[go_memstats_frees_total#LABELS#] \u003C 50000",
            "scope": "",
            "scopeTemplate": ""
          }
        ]
      }
    }
  ]
}
```

###### GET api/v1/sla/:id

e.g. http://sla-manager:8081/api/v1/slas/ExampleApplication_01-Enw6R5Pni7eanXVHtEM8sR

```json
{
  "Message": "Object found",
  "Method": "GetSLA",
  "Resp": "ok",
  "Response": {
    "id": "ExampleApplication_01-Enw6R5Pni7eanXVHtEM8sR",
    "name": "ExampleApplication_01",
    "state": "started",
    "assessment": {
      "total_executions": 14837,
      "total_violations": 14837,
      "x": 2,
      "x_assessment_broken_count": 14837,
      "y": 2,
      "z": 5,
      "level": "Critical",
      "violated": true,
      "first_execution": "2025-05-07T10:15:41.741141409Z",
      "last_execution": "2025-05-12T13:53:41.74250365Z",
      "guarantees": {
        "Processing01": {
          "first_execution": "2025-05-07T10:15:41.741141409Z",
          "last_execution": "2025-05-12T13:53:41.74250365Z",
          "last_values": {
            "go_memstats_frees_total": {
              "key": "go_memstats_frees_total",
              "action": "",
              "namespace": "",
              "value": 239692000,
              "datetime": "2025-05-12T13:53:41.743Z"
            }
          },
          "last_violation": {
            "id": "",
            "agreement_id": "ExampleApplication_01-Enw6R5Pni7eanXVHtEM8sR",
            "guarantee": "Processing01",
            "datetime": "2025-05-12T13:53:41.743Z",
            "constraint": "[go_memstats_frees_total] \u003C 50000",
            "values": [
              {
                "key": "go_memstats_frees_total",
                "action": "",
                "namespace": "",
                "value": 239692000,
                "datetime": "2025-05-12T13:53:41.743Z"
              }
            ],
            "appID": "ExampleApplication_01-Enw6R5Pni7eanXVHtEM8sR"
          }
        }
      }
    },
    "creation": "2025-05-07T10:15:41.314266531Z",
    "expiration": "2026-05-07T10:15:41.31426687Z",
    "details": {
      "guarantees": [
        {
          "name": "Processing01",
          "constraint": "[go_memstats_frees_total] \u003C 50000",
          "query": "[go_memstats_frees_total#LABELS#] \u003C 50000",
          "scope": "",
          "scopeTemplate": ""
        }
      ]
    }
  }
}
```

##### Delete SLA

```bash
curl -k -X DELETE -d @resources/sla_v1.json 'http://<IP>:8081/api/v1/sla/<SLA_ID>'

curl -k -X DELETE -d @resources/sla_v1.json http://sla-manager:8081/api/v1/sla/ExampleApplication_01-Enw6R5Pni7eanXVHtEM8sR
```

----------------------------

## 3. SLAs with scope

The following service definition generates a SLA with a scope value. These SLAs are created with PAUSED state. This means the SLA Manager will wait to have the correspondent value of the scope to start the SLA assessment.

### Service descriptor

Service descriptor is sent to SLA Manager to create the corresponding PAUSED SLA.

**POST api/v1/sla**

```json
{
    "id": {
        "value": "App01"
    },
    "dockerContextDefinitions": [
        {
            "id": "company_premises",
            "imageId": "company_premises:latest"
        }
    ],
    "kpis": [],
    "dockerRoleDefinitions": [
        {
            "id": "Getter",
            "imageId": "getter:latest",
            "hardwareRequirements": [
                "GETTER"
            ],
            "kpis": [{
                "query": "avg_over_time(App01_processing_time[5m]) < 150",
                "scope": "company_premises/building=."
            }]
        }
    ]
}
```

### PAUSED SLA

SLA is created in PAUSED status.

**GET api/v1/slas**

```json
{
  "Message": "Objects found",
  "Method": "GetSLAs",
  "Resp": "ok",
  "Response": [
    {
      "id": "App01-NGXmDaqh8EtTnYmRh2LJkd",
      "name": "App01",
      "state": "paused",
      "assessment": {
        "x": 2,
        "y": 2,
        "z": 5,
        "level": "Unknown",
        "first_execution": "0001-01-01T00:00:00Z",
        "last_execution": "0001-01-01T00:00:00Z"
      },
      "creation": "2025-04-24T11:27:28.2737042+01:00",
      "expiration": "2026-04-24T11:27:28.2737042+01:00",
      "details": {
        "guarantees": [
          {
            "name": "Getter",
            "constraint": "avg_over_time(App01_processing_time[5m]) \u003C 150",
            "query": "[avg_over_time(App01_processing_time#LABELS#[5m])] \u003C 150",
            "scope": "company_premises/building=.",
            "scopeTemplate": "company_premises/building=."
          }
        ]
      }
    }
  ]
}
```

### Send CONTEXT and METRICS to PROMETHEUS

1. Context to Zenoh:

```bash
curl -X PUT -H "content-type:application/json" -d "{\"building\":\"Red\",\"floor\":\"22\",\"room\":\"001\"}" http://zenoh-router:8000/colmena/contexts/ColmenaAgent1/company_premises

curl http://zenoh-router:8000/colmena/contexts/**
```

```
{"key":"colmena/contexts/ColmenaAgent1/company_premises","value":{"building":"Red","floor":"22","room":"001"},"encoding":"application/json","timestamp":""}
```

2. Metrics to Zenoh: 

```bash
curl -X PUT -H "content-type:application/json" -d "{\"company_premises_building\":\"Red\",\"floor\":\"22\",\"room\":\"001\", \"value\":\"24\"}" http://zenoh-router:8000/colmena/metrics/ColmenaAgent1/App01/processing_time
```

3. Check metrics in Prometheus:

```
pp01_processing_time{agent_id="ColmenaAgent1", company_premises_building="Red", floor="22", instance="metrics-etl:8999", job="metrics-etl-colmenagent1", room="001"}	24
```

### STARTED SLA

SLA is updated and set in STARTED status. The SLA Manager can now do the assessment.

**GET api/v1/slas**

```json
{
  "Message": "Objects found",
  "Method": "GetSLAs",
  "Resp": "ok",
  "Response": [
    {
      "id": "App01-NGXmDaqh8EtTnYmRh2LJkd",
      "name": "App01",
      "state": "started",
      "assessment": {
        "x": 2,
        "y": 2,
        "z": 5,
        "level": "Unknown",
        "first_execution": "0001-01-01T00:00:00Z",
        "last_execution": "0001-01-01T00:00:00Z"
      },
      "creation": "2025-04-24T11:27:28.2737042+01:00",
      "expiration": "2026-04-24T11:27:28.2737042+01:00",
      "details": {
        "guarantees": [
          {
            "name": "Getter",
            "constraint": "[avg_over_time(App01_processing_time{company_premises_building=\"Red\"}[5m])] < 150",
            "query": "[avg_over_time(App01_processing_time#LABELS#[5m])] < 150",
            "scope": "company_premises/building=.",
            "scopeTemplate": "company_premises/building=."
          }
        ]
      }
    }
  ]
}
```

#### ASSESSMENT OK

```bash
curl -X PUT -H "content-type:application/json" -d "{\"company_premises_building\":\"Red\",\"floor\":\"3\",\"room\":\"003\", \"value\":\"24\"}" http://zenoh-router:8000/colmena/metrics/ColmenaAgent1/App01/processing_time
```

```json
{
  "Message": "Objects found",
  "Method": "GetSLAs",
  "Resp": "ok",
  "Response": [
    {
      "id": "App01-NGXmDaqh8EtTnYmRh2LJkd",
      "name": "App01",
      "state": "started",
      "assessment": {
        "total_executions": 1,
        "x": 2,
        "y": 2,
        "y_assessment_met_count": 1,
        "z": 5,
        "level": "Met",
        "first_execution": "2025-04-24T11:29:07.2905617+01:00",
        "last_execution": "2025-04-24T11:29:07.2905617+01:00",
        "guarantees": {
          "Getter": {
            "first_execution": "2025-04-24T11:29:07.2905617+01:00",
            "last_execution": "2025-04-24T11:29:07.2905617+01:00",
            "last_values": {
              "avg_over_time(App01_processing_time{company_premises_building=\"Red\"}%5B5m%5D)": {
                "key": "avg_over_time(App01_processing_time{company_premises_building=\"Red\"}%5B5m%5D)",
                "action": "",
                "namespace": "",
                "value": 24,
                "datetime": "2025-04-24T11:29:07.305+01:00"
              }
            }
          }
        }
      },
      "creation": "2025-04-24T11:27:28.2737042+01:00",
      "expiration": "2026-04-24T11:27:28.2737042+01:00",
      "details": {
        "guarantees": [
          {
            "name": "Getter",
            "constraint": "[avg_over_time(App01_processing_time{company_premises_building=\"Red\"}[5m])] \u003C 150",
            "query": "[avg_over_time(App01_processing_time#LABELS#[5m])] \u003C 150",
            "scope": "company_premises/building=.",
            "scopeTemplate": "company_premises/building=."
          }
        ]
      }
    }
  ]
}
```


#### VIOLATION

```bash
curl -X PUT -H "content-type:application/json" -d "{\"company_premises_building\":\"Red\",\"floor\":\"1\",\"room\":\"002\", \"value\":\"124\"}" http://zenoh-router:8000/colmena/metrics/ColmenaAgent1/App01/processing_time
```

----------------------------

## 4. KPI queries

To get the SLAs in a "KPI model", you have to use the following endpoints:

- **GET api/v1/kpis** gets the information about all the SLAs (KPI format)
- **GET api/v1/kpis/:id** gets the information about all the SLAs of a specific service (KPI format)
- **GET api/v1/kpi/:id** gets the information about a specific SLA (KPI format)

Examples:

Use **GET api/v1/kpis** to get all SLAs / KPIs

```json
{
  "Message": "Objects found",
  "Method": "GetKPIs",
  "Resp": "ok",
  "Response": [
    {
      "serviceId": "ExampleApplication_01",
      "slaId": "ExampleApplication_01-f5tjRgFF9HZ5KbgznKamid",
      "KPIs": [
        {
          "roleId": "ExampleApplication_01-f5tjRgFF9HZ5KbgznKamid",
          "query": "[go_memstats_frees_total#LABELS#] \u003C 50000",
          "level": "Broken",
          "value": 0,
          "threshold": "",
          "violations": null,
          "total_violations": 1
        }
      ]
    },
    {
      "serviceId": "ExampleApplication_00",
      "slaId": "ExampleApplication_00-TQG7gzTwZVPjdvimJYrgLQ",
      "KPIs": [
        {
          "roleId": "ExampleApplication_00-TQG7gzTwZVPjdvimJYrgLQ",
          "query": "[go_memstats_frees_total#LABELS#] \u003E 48000",
          "level": "Met",
          "value": 0,
          "threshold": "",
          "violations": null,
          "total_violations": 0
        }
      ]
    }
  ]
}
```

Use **GET api/v1/kpis/:id** to get all SLAs / KPIs from a service

```json
{
  "Message": "Object found",
  "Method": "GetKPIsByServiceId",
  "Resp": "ok",
  "Response": [
    {
      "serviceId": "ExampleApplication_01",
      "slaId": "ExampleApplication_01-f5tjRgFF9HZ5KbgznKamid",
      "KPIs": [
        {
          "roleId": "ExampleApplication_01-f5tjRgFF9HZ5KbgznKamid",
          "query": "[go_memstats_frees_total#LABELS#] \u003C 50000",
          "level": "Critical",
          "value": 0,
          "threshold": "",
          "violations": null,
          "total_violations": 2
        }
      ]
    }
  ]
}
```

Use **GET api/v1/kpi/:id** to get the KPI information of a specific SLA

```json
{
  "Message": "Object found",
  "Method": "GetKPI",
  "Resp": "ok",
  "Response": {
    "serviceId": "ExampleApplication_01",
    "slaId": "ExampleApplication_01-f5tjRgFF9HZ5KbgznKamid",
    "KPIs": [
      {
        "roleId": "ExampleApplication_01-f5tjRgFF9HZ5KbgznKamid",
        "query": "[go_memstats_frees_total#LABELS#] \u003C 50000",
        "level": "Critical",
        "value": 0,
        "threshold": "",
        "violations": null,
        "total_violations": 2
      }
    ]
  }
}
```

----------------------------

## 5. Notifications and violations

Violations and notifications sent to other components (i.e. the endpoint set in **NOTIFICATION_ENDPOINT** environment variable) have the following format:

#### NOTIFICATION

```json
{
    "serviceId": "ExampleApplication_01",
    "slaId": "ExampleApplication_01-Q29AviokdzpGPpdm3WrjD6",
    "KPIs": [{
            "roleId": "ExampleApplication_01-Q29AviokdzpGPpdm3WrjD6",
            "query": "[go_memstats_frees_total#LABELS#] \u003e 50000",
            "level": "Met",
            "value": 58048582,
            "threshold": "",
            "violations": null,
            "total_violations": 0
        }
    ]
}
```

#### VIOLATIONs

Violations are all grouped in a list like the following:

```json
[{
        "serviceId": "ExampleApplication_01",
        "slaId": "ExampleApplication_01-iCP7SemAHCbTXYQcENXjxY",
        "KPIs": [{
                "roleId": "ExampleApplication_01-iCP7SemAHCbTXYQcENXjxY",
                "query": "[go_memstats_frees_total#LABELS#] \u003c 50000",
                "level": "Broken",
                "value": 57198792,
                "threshold": "",
                "violations": [{
                        "id": "",
                        "agreement_id": "ExampleApplication_01-iCP7SemAHCbTXYQcENXjxY",
                        "guarantee": "Processing01",
                        "datetime": "2025-05-21T17:42:38.881+01:00",
                        "constraint": "[go_memstats_frees_total] \u003c 50000",
                        "values": [{
                                "key": "go_memstats_frees_total",
                                "action": "",
                                "namespace": "",
                                "value": 57198792,
                                "datetime": "2025-05-21T17:42:38.881+01:00"
                            }
                        ],
                        "appID": "ExampleApplication_01-iCP7SemAHCbTXYQcENXjxY"
                    }, {
                        "id": "",
                        "agreement_id": "ExampleApplication_01-iCP7SemAHCbTXYQcENXjxY",
                        "guarantee": "Processing01",
                        "datetime": "2025-05-21T17:42:38.881+01:00",
                        "constraint": "[go_memstats_frees_total] \u003c 50000",
                        "values": [{
                                "key": "go_memstats_frees_total",
                                "action": "",
                                "namespace": "",
                                "value": 58516981,
                                "datetime": "2025-05-21T17:42:38.881+01:00"
                            }
                        ],
                        "appID": "ExampleApplication_01-iCP7SemAHCbTXYQcENXjxY"
                    }
                ],
                "total_violations": 1
            }
        ]
    }
]
```

----------------------------

## LICENSES
SLA Manager component is licensed under [Apache License, version 2](LICENSE).