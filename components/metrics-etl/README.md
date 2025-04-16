# Metrics ETL (Zenoh Subscriber & Prometheus Exporter)

**COLMENA Metrics ETL** is a Python-based application that subscribes to **Zenoh** topics, processes the received metrics, and exposes them as **Prometheus** metrics. This enables seamless integration between Zenoh and Prometheus for monitoring and observability.

## Table of Contents

- [Project Structure](#project-structure)
- [Installation](#installation)
- [API Usage](#api-usage)
- [Swagger Documentation](#swagger-documentation)
- [Running Test](#running-test)
- [License](#license)

## Project Structure

components/metrics-etl/
│── src/                       # Source code
│   ├── etl.py                 # Core ETL logic
│   ├── config.py              # Loads configurations
│   ├── cli.py                 # Entry-point
│── config/
│   ├── zenoh_config.json5     # Zenoh configuration file
├── install
│   └── Dockerfile             # Docker setup
│── tests/                     # Unit tests
├── pytest.ini                 # Pytest configuration
├── requirements.txt           # Python dependencies
├── setup.py                   # Installation script for the package
│── README.md                  # Documentation
├── LICENSE                    # Project license

## Installation

### Prerequisites

- Python 3.10+
- pip
- [Zenoh-Python library](https://github.com/eclipse-zenoh/zenoh-python)
- [Prometheus client library](https://github.com/prometheus/client_python)

### Installing with setup.py

Clone the repository:

```bash
git clone [GITHUB-REPO]
cd metrics-etl
```

To install the package:

```bash
pip install .
```

For development (editable mode):

```bash
pip install -e .
```

This allows importing the package in other modules using:

```bash
import metrics_etl
```

### Environment Variables

| Variable            | Default Value               | Description                     |
|---------------------|-----------------------------|---------------------------------|
| `AGENT_ID`          | `COLMENA_AGENT`             | Unique identifier for the agent |
| `ZENOH_CONFIG_FILE` | `config/zenoh_config.json5` | Path to Zenoh configuration     |
| `VERSION`           | `develop`                   | Version identifier              |

### ZENOH Configuration

The application reads its configuration from the ZENOH_CONFIG_FILE environment variable, defaulting to [**config/zenoh_config.json5**](components/metrics-etl/config/zenoh_config.json5).

### Deployment

Build the Docker image:

```bash
docker build -t <registry>/metrics-etl:<version> -f components/metrics-etl/install/Dockerfile .
```

Run the container:

```bash
docker compose -f deploy/compose/docker-compose.yaml up -d metrics-etl
```

## API Usage

The Prometheus metrics will be available at: <http://localhost:8999>

## Swagger Documentation

Available at run time in <http://localhost:8999/docs>

## Running Test

To run unit tests:

```sh
pytest >> report.txt
pylint app --exit-zero --reports y >> qa_report.txt
```

## Contribution

Tech:

- **Tech Stack**: Python, FastAPI
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
