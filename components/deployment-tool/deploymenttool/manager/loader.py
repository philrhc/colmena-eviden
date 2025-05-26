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

# Import system modules
import os
import json
import logging

logger = logging.getLogger(__name__)

class DockerDefinitionLoader:
    """
    Loads and processes Docker definitions from a service's service_description.json file.
    """

    def load_docker_definitions(self, build_path: str, service: str, username: str) -> tuple:
        """
        Loads Docker definitions from a service's service_description.json file.

        This function extracts Docker role and context definitions, updates image IDs 
        with the provided username, and identifies Docker images associated with 
        service artifacts.

        Args:
            build_path (str): The root path where the service build artifacts are located.
            service (str): The name of the service being processed.
            username (str): The username used to prefix Docker image IDs.

        Returns:
            tuple: A dictionary containing the parsed service description and a list of tuples 
                with image IDs and their corresponding paths.

        Raises:
            ValueError: If the service description file is missing, contains invalid JSON, 
                        or any other unexpected error occurs.
        """
        try:
            # Define the base build directory and service description file path
            build_root = os.path.join(build_path, service, "build")
            service_description_path = os.path.join(build_root, "service_description.json")

            # Load the service description file
            with open(service_description_path, "r") as f:
                service_description = json.load(f)

            # Extract Docker role and context definitions
            docker_definitions = {}
            definitions = service_description.get("dockerRoleDefinitions", []) + \
                        service_description.get("dockerContextDefinitions", [])

            # Update image IDs with the provided username
            for definition in definitions:
                image_id = f"{username}/{definition['imageId']}".lower()
                definition["imageId"] = image_id
                docker_definitions[definition["id"]] = image_id

            images = []
            # Iterate over service artifacts inside the build directory
            for artifact in os.scandir(build_root):
                if not artifact.is_dir():
                    continue  # Skip non-directory files

                artifact_path = artifact.path
                sub_artifacts = [artifact_path]

                # If the directory is "context", iterate through subdirectories
                if artifact.name == "context":
                    sub_artifacts = [os.path.join(artifact_path, sub) for sub in os.listdir(artifact_path)]

                # Process each artifact or subdirectory
                for path in sub_artifacts:
                    artifact_name = os.path.basename(path)

                    # Check if Dockerfile exists and if the artifact is in dockerDefinitions
                    dockerfile_path = os.path.join(path, "Dockerfile")
                    if os.path.exists(dockerfile_path) and artifact_name in docker_definitions:
                        images.append((docker_definitions[artifact_name], path))

            return service_description, images

        except FileNotFoundError:
            logger.error(f"Service description file not found at path: {service_description_path}")
            raise ValueError(f"Service description file not found for {service}")
        except json.JSONDecodeError:
            logger.error(f"Invalid JSON format in service description: {service}")
            raise ValueError(f"Invalid JSON format in service description for {service}")
        except Exception as e:
            logger.error(f"Unexpected error loading service {service}: {str(e)}")
            raise ValueError(f"Unexpected error loading service {service}: {str(e)}")