# AGENT Installation

The AGENT is composed by the following components:

- [sla manager](components/sla-manager/README.md)
  - SLA management and assessment golang microservice
  - Prometheus
  - node_exporter
  - Grafana (optional)

- [context awareness manager](components/context-awareness-manager-svc/README.md)
  - Context manager golang microservice

- Other components / tools:
  - [zenoh-prometheus-connector](components/zenoh-prometheus-connector/README.md) python microservice

## Docker images

### Requirements

- git (to download repository)
- Docker and Docker-compose (to create and launch the applications)

### Create Images

- Docker images can be created using the **Makefile**. This file can be edited to change names and versions.

```bash
make build-sla-manager-image

make build-zenoh-prometheus-connector-image
```

```bash
$ docker images
REPOSITORY                   TAG          IMAGE ID       CREATED         SIZE
zenoh-prometheus-connector   1.0          563dd9032098   14 hours ago    154MB
sla-manager                  1.0          06b6f7ee665f   4 days ago      32.2MB
```

These commands create the following images:

- **zenoh-prometheus-connector:1.0**
- **sla-manager:1.0**

#### Launch agent

---------------------------------------------------------
