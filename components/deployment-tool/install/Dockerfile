ARG IMAGE_NAME=python:3.10-slim

# Use an official Python runtime as a parent image
FROM ${IMAGE_NAME}

# Install necessary packages, including Docker CLI
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    curl \
    docker.io && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

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