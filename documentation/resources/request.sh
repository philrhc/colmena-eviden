#!/bin/bash

curl -X POST \
  http://localhost:8080/context \
  -H 'Content-Type: application/json' \
  -d '{
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
