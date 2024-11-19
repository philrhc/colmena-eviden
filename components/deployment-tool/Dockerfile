# Use an official Python runtime as a parent image
FROM python:3.10-slim

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
RUN pip install --no-cache-dir -r requirements.txt

# Expose the port that FastAPI will run on
EXPOSE 8000

# Command to run the FastAPI app with uvicorn
CMD ["uvicorn", "app.rest_app:app", "--host", "0.0.0.0", "--port", "8000"]