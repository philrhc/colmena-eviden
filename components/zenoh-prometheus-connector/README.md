# ZENOH-SUSCRIBER / PROMETHEUS-EXPORTER

The Zenoh-Prometheus-Connector is a Python application responsible for connecting to [**Zenoh**](https://github.com/eclipse-zenoh/zenoh-python) to get the metrics sent to this application, and it is also responsible for sending these metrics to **Prometheus**.

## Installation

### Deploy a Container

1. Clone this repository

    ```sh
    git clone https://github.com/eviden-colmena/colmena-eviden.git
    cd components/zenoh-prometheus-connector
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
    docker build -t <GITHUB_REPO>/zenoh-prometheus-connector -f components/zenoh-prometheus-connector/install/Dockerfile .
    ```

4. Run the Docker container

    ```sh
    docker compose -f install/compose/docker-compose.yaml up -d zenoh-prometheus-connector
    ```

5. Verify Prometheus Metrics

The application exposes metrics on port 8999.

Open your browser and go to <http://localhost:8999/metrics>

## API Usage

SLA Metrics Listener

- Subscribes to tests/**

- Updates Prometheus colmena_sla_metric

Context Metrics Listener

- Subscribes to dockerContextDefinitions/**

- Updates Prometheus colmena_context_metric

### TESTS - PUT METRICS

```bash
curl -X PUT -H "content-type:application/json" -d "2" http://192.168.137.47:8000/tests/planta01/habitacion01

curl -X PUT -H "content-type:application/json" -d "3" http://192.168.137.47:8000/tests/planta01/habitacion02

curl -X PUT -H "content-type:application/json" -d "1" http://192.168.137.47:8000/tests/planta02/habitacion01

curl -X PUT -H "content-type:application/json" -d "2" http://192.168.137.47:8000/tests/planta02/habitacion02

curl -X PUT -H "content-type:application/json" -d "3" http://192.168.137.47:8000/tests/planta02/habitacion03
```

----------------------------------------------

### TESTS: READ METRICS IN PROMETHEUS

```text
colmena_total_people{metric_name="tests"}

colmena_total_people{metric_name="tests", label1="planta01"}

colmena_total_people{metric_name="tests", label1="planta02"}

sum by (metric_name, label1) (colmena_total_people{metric_name="tests", label1="planta01"})
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
