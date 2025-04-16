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
from fastapi import FastAPI
import uvicorn
# Import internal modules
from deploymenttool.api.routes import router
from deploymenttool.api.config import setup_logging

# Import the logging configuration
setup_logging()

# Initialize FastAPI app instance
app = FastAPI(title="COLMENA Deployment Tool", version="1.0")

# Include the router with all the endpoints
app.include_router(router)

def serve():
    uvicorn.run("deploymenttool.main:app", host="0.0.0.0", port=8000)