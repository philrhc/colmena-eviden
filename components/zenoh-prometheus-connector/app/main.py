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
logger = logging.getLogger("COLMENA ETL")

# Retrieve agent ID from environment variable
AGENT_ID = os.getenv("AGENT_ID")

# Define Prometheus Gauges
context_metric = Gauge(f'{AGENT_ID}/colmena_context_metric', 'Context Metric description', ['building', 'floor', 'room'])
role_metric = Gauge(f'{AGENT_ID}/colmena_role_metric', 'Role Metric description', ['building', 'floor', 'room'])

def role_listener(data):
    """
    Process data received for the Role key expression.
    This function processes the data and updates the corresponding Prometheus metric.
    """
    logger.info(">>>> Processing COLMENA - Role metrics ...")
    logger.info(f"Processing Role metric: {data}")

    try:
        # Extract key and value from the data
        key_expr = data.key_expr
        raw_payload = data.payload.to_string()

        # Parse the payload as a dictionary
        payload_dict = json.loads(raw_payload)

        # Extract context values
        building = payload_dict.get("building", "unknown")
        floor = payload_dict.get("floor", "unknown")
        room = payload_dict.get("room", "unknown")
        metric = payload_dict.get("value", "unknown")

        # Update Prometheus Role metric
        role_metric.labels(
            building=building,
            floor=floor,
            room=room,
        ).set(metric)

        logger.info(f"Metric processed: [name={key_expr}, context={building}/{floor}/{room}, value={metric}]")

    except Exception as error:
        logger.error(f"Error processing Role metric: {error}")

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

        # Extract context values
        building = payload_dict.get("building", "unknown")
        floor = payload_dict.get("floor", "unknown")
        room = payload_dict.get("room", "unknown")

        # Update Prometheus Context metric
        context_metric.labels(
            building=building,
            floor=floor,
            room=room,
        ).set(1)
            
        logger.info(f"Metric processed: [name={key_expr}, value={building}/{floor}/{room}]")

    except json.JSONDecodeError as e:
        logger.error(f"Failed to parse payload as JSON: {e}")
    except Exception as error:
        logger.error(f"Error processing context data: {error}")

def main():
    """
    Main function to initialize Zenoh, set up subscriptions, and handle Prometheus metrics.
    """
    logger.info(f"Starting [COLMENA ETL, version {os.getenv("VERSION", "develop")}] ...")

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
        role_subscriber = session.declare_subscriber(f"colmena/metrics/{AGENT_ID}/**", role_listener)
        logger.info(f"Subscribed to colmena/metrics/{AGENT_ID}/**")

        context_subscriber = session.declare_subscriber(f"colmena/contexts/{AGENT_ID}/**", context_listener)
        logger.info(f"Subscribed to colmena/contexts/{AGENT_ID}/**")

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