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

# Import necessary modules
import logging
import os

# Set up logging configuration
def setup_logging():
    """
    Set up logging configuration for the application.
    """
    # Set up logging
    logging.basicConfig(level=logging.DEBUG)
    # Suppress debug logs for Docker SDK and urllib3
    logging.getLogger("docker").setLevel(logging.INFO)
    logging.getLogger("urllib3").setLevel(logging.INFO)

# Environment variables
ZENOH_CONFIG = os.getenv('ZENOH_CONFIG_FILE')
LOGS_DIR = os.getenv('LOGS_DIR', './logs')
