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

import asyncio
import logging
import subprocess
import os
import docker
from docker.errors import DockerException

logger = logging.getLogger(__name__)

class ImageBuilder:
    """
    Builds and pushes Docker images using Docker Buildx with multi-arch support.
    """
    def __init__(self, logs_dir: str, log_url_base: str):
        """
        Initializes the ImageBuilder instance.
        """

        # Set up the logs directory and URL base
        self.logs_dir = logs_dir
        self.log_url_base = log_url_base
        os.makedirs(logs_dir, exist_ok=True)

        try:
            self.client = docker.DockerClient(base_url='unix://var/run/docker.sock', tls=False)
            # Check if Docker daemon is running
            self.client.ping()
            logger.info("Connected to Docker daemon successfully.")
        except DockerException as e:
            logger.error(f"Failed to connect to the Docker daemon: {str(e)}")
            raise DockerException(f"Failed to connect to the Docker daemon: {str(e)}")
        except Exception as e:
            logger.error(f"Unexpected error connecting to Docker: {str(e)}")
            raise Exception(f"Unexpected error connecting to Docker: {str(e)}")

    async def process_images(self, images: list, registry: str, platform: str) -> dict:
        """
        Builds and pushes Docker images for each service role or context.

        This function processes a list of images, building each one from its corresponding
        directory and pushing it to a specified Docker registry. Any logs encountered 
        during the process are logged and returned.

        Args:
            images (list): List of tuples [(image_name, dockerfile_path)]
            registry (str): Docker registry.
            platform (str): Target platforms (comma-separated).

        Returns:
            dict: A dictionary mapping failed image names to their respective error messages.
        """

        logs = {}
        for image_name, dockerfile_path in images:
            try:
                # Tag the image with the appropriate registry
                full_tag = self.image_tagging(image_name, registry)

                # Build and push the image using Docker Buildx for multi-arch support
                status, log_message = await self.build_and_push_image(dockerfile_path, full_tag, platform)
                logs[full_tag] = {"success": status, "log_message": log_message}

            except Exception as e:
                logger.error(f"Unexpected error while processing {image_name}: {e}")
                logs[image_name] = {"success": False, "log_message": str(e)}
                continue

        return logs
    
    def image_tagging(self, image: str, registry: str) -> str:
        """
        Tags image with registry (adds :latest if needed).

        Returns:
            str: Fully qualified image tag.
        """

        if ":" not in image:
            image += ":latest"
        return image if registry == "docker.io" else f"{registry}/{image}"

    async def build_and_push_image(self, path: str, tag: str, platform: str) -> tuple[bool, str]:
            """
            Builds and pushes a Docker image using Docker Buildx.

            This function builds a Docker image from the specified path and tags it 
            with support for multi-architecture builds.
            It captures the build output in a log file and returns the status of the build.
            The log file is stored in the logs directory with a name based on the image tag.

            Args:
                path (str): The path to the directory containing the Dockerfile.
                tag (str): The tag to assign to the built image.
                platform (str): The target platform for the build, specified as a comma-separated string.
            Returns:
                tuple: (True, "") if success; (False, error_message) otherwise.
            Raises:
                ValueError: If the platform format is invalid.
                FileNotFoundError: If the Docker executable is not found.
                subprocess.CalledProcessError: If the build process fails.
                Exception: For any other unexpected errors.
            """
            logging.info(f"Starting build for image {tag} on platforms: {platform}")

            log_filename = f"{tag.split('/')[-1].replace(':', '_')}_build.log"
            log_path = os.path.join(self.logs_dir, log_filename)
            log_url = f"{self.log_url_base}/{log_filename}"

            logger.info(f"Building image {tag} for platforms: {platform}")
            try:
                cmd = [
                    "docker", "buildx", "build",
                    "--debug",
                    "--push",
                    "--platform", platform,
                    "--tag", tag,
                    "--no-cache",
                    path,

                ]

                process = await asyncio.create_subprocess_exec(
                    *cmd,
                    stdout=asyncio.subprocess.PIPE,
                    stderr=asyncio.subprocess.STDOUT
                )

                # Use subprocess.Popen to capture real-time output
                with open(log_path, "w") as log_file:
                    async for line in process.stdout:
                        log_file.write(line.decode('utf-8'))

                await process.wait()

                # Check the return code of the process
                if process.returncode != 0:
                    logging.error(f"Build failed for {tag}, logs at: {log_url}")
                    return False, log_url

                logging.info(f"Image {tag} built successfully for platforms: {platform}")
                logging.info(f"Building logs at: {log_url}")
                return True, log_url

            except (ValueError, FileNotFoundError, subprocess.CalledProcessError) as e:
                logger.error(f"Build error for {tag}: {e}")
                return False, str(e)

            except Exception as e:
                logger.error(f"Unexpected error for image {tag}: {e}")
                return False, str(e)
