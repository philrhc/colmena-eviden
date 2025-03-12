#!/usr/bin/python3
# COLMENA-DESCRIPTION-SERVICE
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

from . import models
from fastapi import FastAPI, Body
from fastapi.responses import JSONResponse
from .rest_manager import RestManager
import logging
import uvicorn

# Set up logging
logging.basicConfig(level=logging.DEBUG)
# Suppress debug logs for Docker SDK and urllib3
logging.getLogger("docker").setLevel(logging.INFO)
logging.getLogger("urllib3").setLevel(logging.INFO)
logger = logging.getLogger(__name__)

# Initialize FastAPI app and RestManager
app = FastAPI(title="COLMENA Deployment Tool", version="1.0")
manager = RestManager()

@app.post("/deployservice/{service}")
async def deployService(service: str, request: models.DeployRequest = Body(...)):
    """
    Deploys a COLMENA Service.

    This endpoint:
    1. Loads Docker definitions from the service's build path.
    2. Builds and pushes role/context Docker images to the specified registry.
    3. Publishes the service description.
    
    Args:
        service (str): Name of the service to deploy.
        request (models.DeployRequest): Request body containing:
            - build_path (str): Path to the service's build directory.
            - username (str): Username with account in the Docker registry.
            - registry (str): Docker registry for pushing images.
            - platform (str): Target platform for multi-arch builds.

    Returns:
        JSONResponse: Status and message about deployment success or failure.
    """
    try:
        # Load docker definitions (service description & image build paths)
        service_definition, images = manager.load_docker_definitions(
            request.build_path, service, request.username
        )

        # Build and push images from docker definitions
        failures = manager.build_and_push_images(images, request.registry, request.platform)
        
        # Publish service description regardless of failures
        if not manager.publish_service_description(service_definition):
            return JSONResponse(status_code=500, content={"detail": "Failed to publish service description."})
        
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
            content={"message": "Images processed and service description published."}
        )
        
    except ValueError as error:
        return JSONResponse(status_code=400, content={"detail": error})
    except Exception as error:
        logger.error(f"Unexpected error during deployment: {error}")
        return JSONResponse(status_code=500, content={"detail": "Internal server error."})
    
def serve():
    uvicorn.run("app.rest_app:app", host="0.0.0.0", port=8000)