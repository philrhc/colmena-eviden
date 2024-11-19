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
from fastapi import HTTPException
import zenoh
import logging

# Set up logging
logging.basicConfig(level=logging.DEBUG)
# Set logging level for Docker SDK and urllib3 to suppress debug logs
logging.getLogger("docker").setLevel(logging.INFO)
logging.getLogger("urllib3").setLevel(logging.INFO)

class RestManager:
    def __init__(self, base_directory: str):
        # Read the service description JSON
        config_path = os.path.join(base_directory, "service_description.json")
        with open(config_path, "r") as f:
            self.service_definition = json.load(f)
        
        self.dockerDefinitions = {}  # Dictionary to store image tags for each artifact
        # Parse roles/context and their corresponding imageIds
        for definition in self.service_definition["dockerRoleDefinitions"] + \
        self.service_definition["dockerContextDefinitions"]:
            self.dockerDefinitions[definition["id"]] = definition["imageId"]

        self.base_directory = base_directory
        self.client = docker.DockerClient(
            base_url='unix://var/run/docker.sock',
            tls=False
        )
        # Load configuration file path from the environment variable
        self.zenoh_config = os.getenv('ZENOH_CONFIG_FILE', 'config/zenoh_config.json5')

    # def build_and_push_images(self):
    #     """
    #     Builds Docker images for each role/context found in the base directory
    #     and pushes Docker images to the specified Docker Hub repository.
    #     """
    #     for artifact in os.listdir(self.base_directory): # Example: 'Processing' or 'Sensing'
    #         artifact_path = os.path.join(self.base_directory, artifact)
    #         # Check if is a directory and if Dockerfile is in the directory
    #         if os.path.isdir(artifact_path) and "Dockerfile" in os.listdir(artifact_path):  
    #             # Only build images for roles found in the roles list
    #             if artifact in self.dockerDefinitions:
    #                 # Construct image tag based on role
    #                 image_tag = f"{self.dockerDefinitions[artifact].lower()}:latest"
    #                 self.build_image(image_tag, artifact_path)
    #                 self.push_image(image_tag)
    #         elif artifact == "context" and os.path.isdir(artifact_path):
    #             for sub_artifact in os.listdir(artifact_path):
    #                 sub_artifact_path = os.path.join(artifact_path, sub_artifact)
    #                 if os.path.isdir(sub_artifact_path) and "Dockerfile" in os.listdir(sub_artifact_path):
    #                     if sub_artifact in self.dockerDefinitions:
    #                         image_tag = f"{self.dockerDefinitions[sub_artifact].lower()}:latest"
    #                         self.build_image(image_tag, sub_artifact_path)
    #                         self.push_image(image_tag)

    def build_and_push_images(self):
        """
        Builds Docker images for each role/context found in the base directory
        and pushes Docker images to the specified Docker Hub repository.
        """
        failures = {}
        for artifact in os.listdir(self.base_directory):  # Example: 'Processing', 'Sensing', or 'context'
            artifact_path = os.path.join(self.base_directory, artifact)

            # Check that artifact_path is a directory; skip if it's not
            if not os.path.isdir(artifact_path):
                continue  # Skip non-directory files

            # Determine if we are dealing with a role artifact or 'context' directory
            if artifact == "context":
                # Iterate through subdirectories in 'context' directory
                sub_artifacts = [
                    os.path.join(artifact_path, sub_artifact)
                    for sub_artifact in os.listdir(artifact_path)
                    if os.path.isdir(os.path.join(artifact_path, sub_artifact))
                ]
            else:
                # Otherwise, treat the main artifact as a single directory to check
                sub_artifacts = [artifact_path]

            # Process each artifact or subdirectory as needed
            for path in sub_artifacts:
                artifact_name = os.path.basename(path)
                # Check if Dockerfile exists and if the artifact is in dockerDefinitions
                if "Dockerfile" in os.listdir(path) and artifact_name in self.dockerDefinitions:
                    image_tag = f"{self.dockerDefinitions[artifact_name].lower()}:latest"
                    try:
                        # TODO: Docker login for DockerHub registry.
                        self.build_image(image_tag, path)
                        self.push_image(image_tag)
                    except Exception as error:
                        logging.error(str(error))
                        failures[image_tag] = str(error)
                        continue
            
        # If there were any failures, raise an HTTPException with the error details
        if failures:
            raise HTTPException(
                status_code=207,  # Multi-Status: Partial success
                detail=failures  # Returning the failures dictionary as the response
            )

    # def build_image(self, image_tag: str, build_context: str):
    #     """
    #     Builds a Docker image using Buildah for the given directory.
    #     """
    #     buildah_build_command = [
    #         "buildah", "bud",
    #         "-t", f"{self.repo_url}/{image_tag}",
    #         build_context  # Directory containing the Dockerfile
    #     ]
        
    #     logging.info(f"Building image {image_tag}...")
    #     result = subprocess.run(buildah_build_command, capture_output=True, text=True)
    #     if result.returncode != 0:
    #         logging.error(f"Error building image {image_tag}:\n{result.stderr}")
    #         raise HTTPException(status_code=500, detail=f"Failed to build image {image_tag}")

    #     logging.info(f"Image {image_tag} built successfully.")

    # def push_image(self, image_tag: str):
    #     """
    #     Pushes a Docker image to the repository using Buildah.
    #     """
    #     buildah_push_command = [
    #         "buildah", "push", "--tls-verify=false",
    #         f"{self.repo_url}/{image_tag}"
    #     ]

    #     logging.info(f"Pushing image {image_tag} to the registry...")
    #     result = subprocess.run(buildah_push_command, capture_output=True, text=True)
    #     if result.returncode != 0:
    #         logging.error(f"Error pushing image {image_tag}:\n{result.stderr}")
    #         raise HTTPException(status_code=500, detail=f"Failed to push image {image_tag}")

    #     logging.info(f"Image {image_tag} pushed successfully to {self.repo_url}.")

    def build_image(self, image_tag: str, build_context: str):
        """
        Builds a Docker image for the given directory using Docker SDK for Python.
        """
        logging.info(f"Building image {image_tag}...")

        try:
            # Build the Docker image from the specified build context (directory with Dockerfile)
            image, logs = self.client.images.build(path=build_context, tag=image_tag)
            
            # Capture and log the build process
            for log in logs:
                if 'error' in log:
                    raise Exception(f"Failed to build image {image_tag}: {log['error']}")

            logging.info(f"Image {image_tag} built successfully.")

        except docker.errors.BuildError as e:
            raise Exception(f"Failed to build image {image_tag}: {str(e)}")

    def push_image(self, image_tag: str):
        """
        Pushes a Docker image to the repository using subprocess and HTTP.
        """
        logging.info(f"Pushing image {image_tag} to the registry...")
        try:
            # Push the image to the Docker registry
            for line in self.client.images.push(image_tag, stream=True, decode=True):
                if 'error' in line:
                    raise Exception(f"Error pushing image {image_tag}: {line['error']}")
            logging.info(f"Image {image_tag} pushed successfully to repository.")

        except docker.errors.APIError as e:
            raise Exception(f"Error pushing image {image_tag}: {str(e)}")


    def publish_service_description(self):
        """
        Publish service definition to zenoh under a specific keyexpr

        Parameters:
            service_definition: 
        """
        logging.info("Publishing service description...")
        try:
            zenoh_session = zenoh.open(zenoh.Config.from_file(self.zenoh_config))
            payload = json.dumps(self.service_definition)
            zenoh_session.put("colmena_service_definitions", payload.encode('utf-8'))
            logging.info(f"Service description successfully published to Zenoh.")
        except Exception as error:
            logging.error(f"Failed to publsishing service description to Zenoh: {str(error)}")
    