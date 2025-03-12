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
COPY components/zenoh-prometheus-connector .

# Install dependencies
RUN pip install -e .

EXPOSE 8999

CMD ["python", "app/main.py"]