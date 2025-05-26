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
import zenoh
import os
import json

logger = logging.getLogger(__name__)

class Publisher:
    """
    A class to manage the publishing of service descriptions to Zenoh.
    This class provides methods to publish service descriptions and handle Zenoh connections.
    """
    def __init__(self, zenoh_config: str = "config/zenoh_config.json5"):
        """
        Initializes the Publisher instance.

        Args:
            zenoh_config (str): Path to the Zenoh configuration file.
        """
        self.zenoh_config = zenoh_config
        self.zenoh_session = None

        # Check if the Zenoh configuration file exists
        if not os.path.isfile(self.zenoh_config):
            logger.error(f"Zenoh configuration file not found: {self.zenoh_config}")
            raise FileNotFoundError(f"Zenoh configuration file not found: {self.zenoh_config}")

        # Initialize Zenoh session
        try:
            self.zenoh_session = zenoh.open(zenoh.Config.from_file(self.zenoh_config))
            logger.info("Zenoh session initialized successfully.")
        except zenoh.ZError as error:
            logger.error(f"Failed to initialize Zenoh session: {str(error)}")
            raise zenoh.ZError(f"Failed to initialize Zenoh session: {str(error)}")
    
    def publish_service_description(self, service_description: dict) -> bool:
        """
        Publishes the service description to Zenoh.

        Args:
            service_description (dict): The service description to publish.
        Returns:
            bool: True if the publication was successful, False otherwise.
        Raises:
            zenoh.ZError: If there is an error connecting to Zenoh or publishing the message.
        """
        logging.info("Publishing service description to Zenoh...")
        try:
            # Convert the service description to bytes
            service_description = json.dumps(service_description).encode('utf-8')
            # Publish the service description to Zenoh
            self.zenoh_session.put("colmena_service_descriptions", service_description)
            logger.info("Service description published successfully.")
            return True
        except zenoh.ZError as error:
            logger.error(f"Failed to publish service description: {str(error)}")
            return False