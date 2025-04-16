# Use an official Python runtime as a parent image
FROM registry.atosresearch.eu:18520/library/python:3.10-slim

# Install necessary packages, including Docker CLI
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    curl \
    docker.io && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

ENV DOCKER_BUILDKIT=1
ENV BUILDKIT_PROGRESS=plain
ENV DOCKER_CLI_EXPERIMENTAL=enabled

ARG BUILDX_URL=https://github.com/docker/buildx/releases/latest/download/buildx-v0.22.0.linux-amd64

RUN mkdir -p /usr/local/lib/docker/cli-plugins && \
    curl -L ${BUILDX_URL} \
    -o /usr/local/lib/docker/cli-plugins/docker-buildx && \
    chmod +x /usr/local/lib/docker/cli-plugins/docker-buildx

# Set the working directory in the container
WORKDIR /app

# Copy the application directory into the container
COPY components/deployment-tool .

# Install dependencies
RUN pip install -e .

# Expose the port that FastAPI will run on
EXPOSE 8000

# Command to run the FastAPI app with uvicorn via serve in rest_app
CMD ["run_app"]