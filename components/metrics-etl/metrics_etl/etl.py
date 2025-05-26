"""
Copyright © 2024 EVIDEN

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

This work has been implemented within the context of COLMENA project.
"""
import logging
import json
from prometheus_client import Gauge, start_http_server
from zenoh import open, Config, Sample
from metrics_etl.config import AGENT_ID, ZENOH_CONFIG_FILE, PROMETHEUS_PORT, VERSION

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger("COLMENA Metrics ETL")

# Dictionary to store dynamically created Gauges
metrics_registry = {}

def get_or_create_gauge(name, documentation, labels):
    """
    Creates or retrieves a Gauge metric with dynamic labels.
    """
    label_tuple = tuple(sorted(labels)) # Sort labels for consistent lookup

    if (name, label_tuple) not in metrics_registry:
        metrics_registry[(name, label_tuple)] = Gauge(
            name=name,
            documentation=documentation,
            labelnames=label_tuple
        )

    return metrics_registry[(name, label_tuple)]

# Define listeners
def role_listener(data: Sample):
    """
    Processes data for the Role metric and updates Prometheus.
    """
    logger.info(f"Processing COLMENA - Role metrics ...")

    try:
        # Parse the payload
        key_expr = str(data.key_expr)
        payload_dict = json.loads(data.payload.to_string())

        # Extract metric name and service name from the key expression
        topic_parts = key_expr.split("/")
        metric_name = topic_parts[-1]
        service_name = topic_parts[-2]
        metric = f"{service_name}_{metric_name}"
        metric_doc = f"Role Metric for {metric}"

        # Extract metric value and remove it from labels
        metric_value = payload_dict.pop("value", 0)

        # Labels
        labels = {k.replace("/", "_"): v for k, v in payload_dict.items()}

        # Get or create a Gauge with the appropriate labels
        gauge = get_or_create_gauge(metric, metric_doc, labels.keys())

        # Update Prometheus metric with dynamic labels
        gauge.labels(
            **labels,
        ).set(metric_value)

        logger.info(f"Processed {metric}: [labels={labels}, value={metric_value}]")

    except json.JSONDecodeError as e:
        logger.error(f"JSON parsing error: {e}")
    except Exception as error:
        logger.error(f"Error processing Role metric: {error}")

def context_listener(data: Sample):
    """Processes data for the Context metric and updates Prometheus."""
    logger.info("Processing COLMENA - Context metrics ...")

    try:
        # Parse the payload
        key_expr = str(data.key_expr)
        payload_dict = json.loads(data.payload.to_string())

        # Extract metric name and service name from the key expression
        topic_parts = key_expr.split("/")
        context = topic_parts[-1]
        metric = "colmena_context_metric"
        context_doc = f"Context Metric for {context}"

        # Labels
        labels = {
            "context": context,
            **payload_dict,
        }

        # Get or create a Gauge with the appropriate labels
        gauge = get_or_create_gauge(metric, context_doc, labels.keys())

        # Update Prometheus metric with dynamic labels
        gauge.labels(
            **labels,
        ).set(1)

        logger.info(f"Processed {context}: [labels={labels}]")

    except json.JSONDecodeError as e:
        logger.error(f"JSON parsing error: {e}")
    except Exception as error:
        logger.error(f"Error processing Context metric: {error}")

# Start the ETL process
def start_etl():
    """Starts the ETL process with Zenoh and Prometheus."""
    logger.info(f"Starting Metrics ETL (Version: {VERSION}) on Agent: {AGENT_ID}")

    start_http_server(PROMETHEUS_PORT)
    logger.info(f"Prometheus metrics server running on port {PROMETHEUS_PORT}")

    try:
        # Initialize Zenoh session
        session = open(Config.from_file(ZENOH_CONFIG_FILE))
        logger.info(f"Using Zenoh configuration file: {ZENOH_CONFIG_FILE}")

        # Declare subscribers
        role_subscriber = session.declare_subscriber(f"colmena/metrics/{AGENT_ID}/**", role_listener)
        logger.info(f"Subscribed to colmena/metrics/{AGENT_ID}/**")

        context_subscriber = session.declare_subscriber(f"colmena/contexts/{AGENT_ID}/**", context_listener)
        logger.info(f"Subscribed to colmena/contexts/{AGENT_ID}/**")

        while True:
            pass  # Keeps running

    except Exception as error:
        logger.error(f"Failed to start ETL: {error}")
    finally:
        session.close()
        logger.info("Zenoh session closed.")