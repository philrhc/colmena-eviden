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
from fastapi import APIRouter, Body
from fastapi.responses import JSONResponse
from fastapi.responses import FileResponse
# Import internal modules
from deploymenttool.api import models
from deploymenttool.manager.buildx import BuildxManager
from deploymenttool.manager.images import ImageBuilder
from deploymenttool.manager.loader import DockerDefinitionLoader
from deploymenttool.manager.publisher import Publisher
from deploymenttool.api.config import LOGS_DIR, ZENOH_CONFIG
import deploymenttool.manager.utils as utils
# Import system modules
import logging
import os

logger = logging.getLogger(__name__)

# Initialize router
router = APIRouter()

# Initialize the Docker Builder Buildx
buildx = BuildxManager()
# Initialize the Docker manager
builder = ImageBuilder(logs_dir=LOGS_DIR, log_url_base="/logs")
# Initialize the Docker definition loader
loader = DockerDefinitionLoader()
# Initialize the Zenoh publisher
publisher = Publisher(zenoh_config=ZENOH_CONFIG)

@router.get("/")
async def root():
    """
    Root endpoint to check if the service is running.
    """
    return JSONResponse(status_code=200, content={"message": "COLMENA Deployment Tool is running."})

@router.get("/health")
async def health_check():
    """
    Health check endpoint to verify the service's status.
    """
    return JSONResponse(status_code=200, content={"message": "Service is healthy."})

@router.get("/logs/{filename}")
def get_log(filename: str):
    """
    Endpoint to retrieve a log file.
    """
    log_path = os.path.join(LOGS_DIR, filename)
    if os.path.exists(log_path):
        return FileResponse(log_path)
    return JSONResponse(status_code=404, content={"error": "Log not found"})

@router.post("/deploy-service/{service_name}")
async def deploy_service(service_name: str, request: models.DeployRequest = Body(...)):
    """
    Endpoint to deploy a service by building and pushing Docker images.
    
    The deployment process includes the following steps:
    1. Load Docker definitions from the service's build path.
    2. Build and push role/context Docker images to the specified registry.
    3. Publish the service description.
    4. Return a status message indicating success or failure.

    The request body should contain:
    - build_path (str): Path to the service's build directory.
    - username (str): Username with account in the Docker registry.
    - registry (str): Docker registry for pushing images.
    - platform (str): Target platform for multi-arch builds.

    The response will include:
    - status (int): HTTP status code indicating the result of the operation.
    - message (str): A message indicating the success or failure of the deployment.
    - failures (list): A list of failed images with error details, if any.
    """

    try:
        # Validate platform format
        utils.validate_platforms(request.platform)

        # Load docker definitions (service description & image build paths)
        service_definition, images = loader.load_docker_definitions(
            request.build_path, service_name, request.username
        )

        # Create docker builder
        if not buildx.create_docker_buildx(platform=request.platform):
            return JSONResponse(status_code=500, content={"detail": "Failed to create Docker Buildx."})
        
        # Build and push images from docker definitions
        logs = await builder.process_images(images, request.registry, request.platform)
        buildx.remove_docker_buildx()
        
        # Publish service description regardless of failures
        if not publisher.publish_service_description(service_definition):
            return JSONResponse(status_code=500, content={"detail": "Failed to publish service description."})
        
        # Check if there are any failures
        failures = {image: data for image, data in logs.items() if not data["success"]}
        
        # Return response based on image build/push results
        if failures:
            return JSONResponse(
                status_code=207,  # Multi-Status: Partial success
                content={
                    "detail": "Some images failed to build/push, but service description was published.",
                    "failures": failures  # List of failed images with error details
                }
            )
        
        # Return success response
        return JSONResponse(
            status_code=200,
            content={
                "detail": "Images processed and service description published.",
                "logs": logs
            }
        )
        
    except ValueError as error:
        return JSONResponse(status_code=400, content={"detail": str(error)})
    except Exception as error:
        logger.error(f"Unexpected error during deployment: {error}")
        return JSONResponse(status_code=500, content={"detail": "Internal server error."})
