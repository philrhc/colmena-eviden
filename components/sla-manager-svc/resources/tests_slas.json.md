## 0. SERVICE DESCRIPTION

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

## 1. PAUSED

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

## 2. STARTED

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

## 3. ASSESSMENT OK

```bash
curl -X PUT -H "content-type:application/json" -d "{\"company_premises_building\":\"Red\",\"floor\":\"122\",\"room\":\"Rest Room\", \"value\":\"424\"}" http://192.168.137.47:8000/colmena/metrics/ColmenaAgent1/App01/processing_time
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

```json
 {"serviceId":"App01","KPIs":[{"roleId":"App01-NGXmDaqh8EtTnYmRh2LJkd","query":"[avg_over_time(App01_processing_time#LABELS#[5m])] \u003c 150","level":"Met","value":null,"threshold":"","violations":null,"total_violations":0}]}

```

## 4. ASSESSMENT KO

```bash
curl -X PUT -H "content-type:application/json" -d "{\"company_premises_building\":\"Red\",\"floor\":\"22\",\"room\":\"Rest Room\", \"value\":\"424\"}" http://192.168.137.47:8000/colmena/metrics/ColmenaAgent1/App01/processing_time
```

```json


```



```json


```