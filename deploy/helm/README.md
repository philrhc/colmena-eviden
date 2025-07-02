# Monitoring Manager

Monitoring Manager: Prometheus compatible long term storage for monitoring multi-cluster telemetry data

# Interfaces documentation

Telemetry data is pushed to the Monitoring Manager Thanos server through the Thanos Receiver following the [Prometheus Remote-Write Specification](https://prometheus.io/docs/concepts/remote_write_spec/) (gRPC)

Monitoring Manager also serves telemetry data through the Thanos Querier under the [Prometheus HTTP v1 API](https://prometheus.io/docs/prometheus/latest/querying/api/)

# Getting Started

Monitoring Manager is deployed as a Helm Chart.

Create a helm `values.yaml` file in the root directory with your customized values (Thanos receive and query hostnames).

Deploy using Helm:
```sh
helm upgrade -i  monitoring-manager monitoring-manager/ -f values.yaml
```