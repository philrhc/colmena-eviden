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
import subprocess

logger = logging.getLogger(__name__)

class BuildxManager:
    """
    Manages Docker Buildx for multi-architecture builds.
    """
    def __init__(self, builder_name="colmena-builder", network_name="compose_colmena-net"):
        """
        Initializes the BuildxManager instance.
        """
        self.builder_name = builder_name
        self.builder_container = f"buildx_buildkit_{builder_name}0"
        self.network_name = network_name

    def create_docker_buildx(self, platform: str) -> bool:
        """
        Creates and configures a Docker Buildx builder instance for multi-architecture builds.

        This includes:
        1. Removing any existing builder.
        2. Validating the target platform string.
        3. Creating the Buildx builder.
        4. Connecting it to a specific Docker network (usually from docker-compose).

        Args:
            platform (str): Comma-separated string of platforms (e.g., 'linux/amd64,linux/arm64').

        Returns:
            bool: True if builder was successfully created and connected, False otherwise.
        """
        
        try:
            # Remove any existing Docker Buildx builder instance
            self.remove_docker_buildx()

            # Create the Docker Buildx builder and set it as the active builder
            logging.info(f"Creating Docker Buildx builder for platforms: {platform}...")
            cmd = [
                "docker", "buildx", "create",
                "--name", self.builder_name,
                "--platform", platform,
                "--use",
                "--bootstrap",
            ]
            subprocess.run(cmd, check=True, capture_output=True, text=True)
            logging.info("Docker Buildx builder created successfully.")
            return True
        
        except subprocess.CalledProcessError as e:
            logging.error(f"Failed to create Docker Buildx builder: {str(e)}")
            return False

    def remove_docker_buildx(self):
        """
        Removes the Docker Buildx builder instance.
        """
        logging.debug("Removing Docker Buildx builder...")
        try:
            cmd = [
                "docker", "buildx", "rm",
                self.builder_name,
            ]
            subprocess.run(cmd, check=True, capture_output=True, text=True)
            logging.debug("Docker Buildx builder removed successfully.")

        except subprocess.CalledProcessError as e:
            if "no builder" not in e.stderr:
                logging.error(f"Failed to remove Docker Buildx builder: {str(e)}")
     