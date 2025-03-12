#!/usr/bin/pythonCOLMENA-DESCRIPTION-SERVICE
# Copyright © 2024 EVIDEN

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

# http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This work has been implemented within the context of COLMENA project.

# -*- coding: utf-8 -*-

import os
import docker
import json
import zenoh
import logging
import subprocess
import re

logger = logging.getLogger(__name__)

class RestManager:
    def __init__(self):
        """
        Initializes the Docker client and performs Docker login if credentials are set.
        """
        self.client = docker.DockerClient(base_url='unix://var/run/docker.sock', tls=False)
        self.zenoh_config = os.getenv('ZENOH_CONFIG_FILE', 'config/zenoh_config.json5')

    def docker_login(self, registry: str, username: str, password: str):
        """Logs into a Docker registry."""
        logger.info(f"Logging into {registry} as {username}")
        try:
            self.client.login(username=username, password=password, registry=registry)
            logger.info("Docker login successful.")
            return True
        except docker.errors.APIError as error:
            logger.error(f"Docker login failed: {str(error)}")
            return False
            
    def load_docker_definitions(self, build_path: str, service: str, username: str):
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
            logger.error(f"Service description file not found: {service}")
            raise ValueError(f"Service description file not found for {service}")
        except json.JSONDecodeError:
            logger.error(f"Invalid JSON format in service description: {service}")
            raise ValueError(f"Invalid JSON format in service description for {service}")
        except Exception as e:
            logger.error(f"Unexpected error loading service {service}: {str(e)}")
            raise ValueError(f"Unexpected error loading service {service}: {str(e)}")

    def build_and_push_images(self, images: list, registry: str, platform: str):
        """
        Builds and pushes Docker images for each service role or context.

        This function processes a list of images, building each one from its corresponding
        directory and pushing it to a specified Docker registry. Any failures encountered 
        during the process are logged and returned.

        Args:
            images (list): A list of tuples, where each tuple contains:
                        - image (str): The base image name.
                        - path (str): The directory containing the Dockerfile.
            registry (str): The target Docker registry for pushing images.
            platform (str): The platform architecture for multi-arch builds (not currently used).

        Returns:
            dict: A dictionary mapping failed image names to their respective error messages.
        """
        failures = {}
        for image, path in images:
            try:
                # Tag the image with the appropriate registry
                tagged_image = self.image_tagging(image, registry)
                # Build the Docker image from the specified path
                self.build_image(path, tagged_image)
                # Push the image to the registry
                self.push_image(registry, tagged_image)

                # Placeholder for multi-architecture builds
                # success, error_message = self.build_and_push_image_multiarch(path, tagged_image, platform)
                # if not success:
                #     failures[tagged_image] = error_message

            except Exception as error:
                logger.error(f"Error processing {image}: {str(error)}")
                failures[image] = str(error)
                continue

        return failures
    
    def image_tagging(self, image: str, registry: str):
        """
        Formats the Docker image name by ensuring it has a proper tag and registry prefix.

        - If no tag is provided, ':latest' is appended.
        - If a Docker registry is specified, it is prefixed to the image name.
        - If the registry is 'docker.io', the image name is returned as-is since Docker handles it.
        """
        if ":" not in image:
            image = f"{image}:latest"  # Default to 'latest' if no tag is provided

        # If repository is docker.io, use only image_name (Docker handles it internally)
        return image if registry == "docker.io" else f"{registry}/{image}"
    
    def build_image(self, path: str, image: str):
        """Builds a Docker image using the Docker SDK."""
        logging.info(f"Building image {image}...")

        try:
            # Build the image with no cache and remove intermediate containers
            image, logs = self.client.images.build(path=path, tag=image, nocache=True, rm=True)
            
            # Capture and log build output
            for log in logs:
                if 'error' in log:
                    raise Exception(f"Failed to build image {image}: {log['error']}")

            logging.info(f"Image {image} built successfully.")
        except docker.errors.BuildError as error:
            raise Exception(f"Failed to build image {image}: {str(error)}")

    def push_image(self, repository: str,  image: str):
        """Pushes a Docker image to the specified registry."""
        logging.info(f"Pushing image {image} to the registry {repository}...")
        try:
            # Push the image and log any errors
            for line in self.client.images.push(image, stream=True, decode=True):
                if 'error' in line:
                    raise Exception(f"Error pushing {image} image: {line['error']}")
                
            logging.info(f"Image {image} pushed successfully to {repository}.")
        except docker.errors.APIError as error:
            raise Exception(f"Error pushing image {image}: {str(error)}")

    def publish_service_description(self, service_description):
        """Publishes the service description to Zenoh."""
        logging.info("Publishing service description to Zenoh...")
        try:
            zenoh_session = zenoh.open(zenoh.Config.from_file(self.zenoh_config))
            payload = json.dumps(service_description)
            zenoh_session.put("colmena_service_descriptions", payload.encode('utf-8'))
            logging.info(f"Service description successfully published to Zenoh.")
            return True
        except zenoh.ZError as error:
            logging.error(f"Failed to publish service description to Zenoh: {str(error)}")
            return False

    def validate_platforms(self, platforms: list):
        """Validates that each platform follows the format os/arch[/variant]."""
        pattern = re.compile(r'^[a-zA-Z0-9]+/[a-zA-Z0-9]+(?:/[a-zA-Z0-9]+)?$')
        for platform in platforms:
            if not pattern.match(platform):
                raise ValueError(f"Invalid platform format: {platform}")
            
    def build_and_push_image_multiarch(self, path: str, tag: str, platform: str):
        """Builds a Docker image using Docker Buildx for multiple architectures.
        Returns True if successful, False if there is any failure along with a detailed error message."""
        logging.info(f"Building image {tag} for platforms: {platform}...")

        try:
            self.validate_platforms(platform)

            cmd = [
                "docker", "buildx", "build",
                "--platform", platform,
                "-t", tag,
                path,
                "--rm=true",
                "--no-chache=true",
                "--push"  # Required to support multi-platform builds
            ]

            subprocess.run(cmd, check=True, capture_output=True, text=True)
            logging.info(f"Image {tag} built successfully for platforms: {platform}")
            return True, ""

        except ValueError as ve:
            error_message = f"Platform validation failed: {ve}"
            logging.error(error_message)
            return False, error_message

        except FileNotFoundError:
            error_message = "Docker is not installed or not in PATH."
            logging.error(error_message)
            return False, error_message

        except subprocess.CalledProcessError as error:
            error_message = f"Failed to build and push image {tag}: {error.stderr}"
            logging.error(error_message)
            return False, error_message

        except Exception as e:
            error_message = f"Unexpected error while building image {tag}: {str(e)}"
            logging.error(error_message)
            return False, error_message