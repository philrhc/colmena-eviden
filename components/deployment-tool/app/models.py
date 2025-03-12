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

from typing import Optional
from pydantic import BaseModel, Field

# Define the request body schema
class DeployRequest(BaseModel):
    build_path: str = Field(..., description="Path to the service's build directory.", examples=["/path/to/build"])
    username: str = Field(..., description="Username with account in the Docker registry.", examples=["user"])
    registry: str = Field(default="docker.io", description="Docker registry", examples=["localhost:5000"])
    platform: Optional[str] = Field(None, description="Target platform for multi-arch builds.", examples=["linux/amd64,linux/arm64"])