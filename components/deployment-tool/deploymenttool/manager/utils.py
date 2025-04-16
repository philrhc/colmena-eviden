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

import logging
import re

logger = logging.getLogger(__name__)

@staticmethod
def validate_platforms(platforms: str):
    """
    Validates that each platform follows the format os/arch[/variant].

    The expected format is:
    - `os`: The operating system (e.g., linux, windows).
    - `arch`: The architecture (e.g., amd64, arm64).
    - `variant` (optional): Additional details about the architecture (e.g., v7 for ARM).

    Args:
        platforms (list): A list of platform strings to validate.
    Raises:
        ValueError: If any platform string does not match the expected format.
    """
    platform_list = platforms.split(',')
    pattern = re.compile(r'^[a-zA-Z0-9]+/[a-zA-Z0-9]+(?:/[a-zA-Z0-9]+)?$')
    for platform in platform_list:
        # Check if the platform string matches the expected format
        if not pattern.match(platform):
            raise ValueError(f"Invalid platform format: {platform}")