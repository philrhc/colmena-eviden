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

from fastapi import FastAPI, HTTPException
from .models import InputData
from .rest_manager import RestManager
import logging

app = FastAPI()

# Set up logging
logging.basicConfig(level=logging.DEBUG)

@app.post("/build-and-push/")
async def build_and_push(input_data: InputData):
    """
    Endpoint to process images: build, push, and publish a service description.
    """
    manager = RestManager(input_data.base_directory)

    try:
        manager.build_and_push_images()
        # TODO: if any image fail to build/push, respond HTTP 207 and dont publish to Zenoh. 
        manager.publish_service_description()
    except HTTPException as error:
        raise error
    except Exception as error:
        raise HTTPException(status_code=500, detail=f"An error occurred: {str(error)}")

    return {"status": "success", "message": "Images processed and service description published successfully."}
