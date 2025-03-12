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
import os, json
from prometheus_client import Gauge, start_http_server
from zenoh import open, Config

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger("Zenoh-Prometheus-Connector")

colmena_total_people = Gauge('colmena_metric1', 'metric description', ['metric', 'path', 'desc'])
# Define Prometheus Gauges
sla_metric = Gauge('colmena_sla_metric', 'SLA Metric description', ['metric', 'path', 'desc'])
context_metric = Gauge('colmena_context_metric', 'metric description', ['metric', 'name', 'value', 'desc'])

def sla_listener(data):
    """
    Process data received for the SLA key expression.
    This function processes the data and updates the corresponding Prometheus metric.
    """
    logger.info(">>>> Processing COLMENA - SLA metrics ...")
    logger.info(f"Processing SLA metric: {data}")

    try:
        # Extract key and value from the data
        key_expr = data.key_expr
        payload = data.payload.to_string()

        # Update Prometheus SLA metric
        sla_metric.labels(
            metric="colmena_sla",
            path=key_expr,
            desc='SLA metric description'
        ).set(payload)

        logger.info(f"Metric processed: [metric=colmena_sla, path={key_expr}]")

    except Exception as error:
        logger.error(f"Error processing SLA metric: {error}")

def context_listener(data):
    """
    Process data received for the Context key expression.
    This function processes the data and updates the corresponding Prometheus metric.
    """
    logger.info(">>>> Processing COLMENA - Context values ...")
    logger.info(f"Processing Context: {data.key_expr}")

    try:
        # Extract key and value from the data
        key_expr = data.key_expr
        raw_payload = data.payload.to_string()
        
        # Parse the payload as a dictionary
        payload_dict = json.loads(raw_payload)

        # Iterate over all fields in the dictionary
        for field, value in payload_dict.items():
            # Create a flexible value representation
            payload_value = f"{field}: {value}"

            # Update Prometheus context metric
            context_metric.labels(
                metric="colmena_context", 
                name=key_expr, 
                value=payload_value, 
                desc=f'Context value for {field}'
            ).set(1)
            
            logger.info(f"Metric processed: [name={key_expr}, value={payload_value}]")

    except json.JSONDecodeError as e:
        logger.error(f"Failed to parse payload as JSON: {e}")
    except Exception as error:
        logger.error(f"Error processing context data: {error}")


def main():
    """
    Main function to initialize Zenoh, set up subscriptions, and handle Prometheus metrics.
    """
    logger.info("Starting [ZENOH-PROMETHEUS-CONNECTOR, version develop] ...")

    # Start the Prometheus HTTP server
    start_http_server(8999)
    logger.info("Prometheus metrics server running on port 8999")

    # Load configuration file path from the environment variable
    zenoh_config = os.getenv('ZENOH_CONFIG_FILE', 'config/zenoh_config.json5')
    logger.info(f"Using Zenoh configuration file: {zenoh_config}")

    try:
        # Initialize Zenoh session
        session = open(Config.from_file(zenoh_config))
        logger.info(f"Zenoh session initialized successfully")

        # Declare subscribers
        session.declare_subscriber("tests/**", sla_listener)
        logger.info("Subscribed to tests/**")

        session.declare_subscriber("dockerContextDefinitions/**", context_listener)
        logger.info("Subscribed to dockerContextDefinitions/**")

    except Exception as error:
        logger.error(f"Failed to subscribe to key expressions: {error}")

    # Keep the application running indefinitely
    logger.info("Zenoh-Prometheus Connector is running...")
    try:
        # Keep the application running indefinitely
        while True:
            pass  # Infinite loop to keep the application alive
    except KeyboardInterrupt:
        logger.info("Shutting down Zenoh-Prometheus Connector...")
    finally:
        session.close()
        logger.info("Zenoh session closed.")

if __name__ == '__main__':
    # Run the main function
    main()