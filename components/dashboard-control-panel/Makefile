.PHONY: all build push login clean

# Nexus registry and image details
registry_host := registry.atosresearch.eu:18512
image_path := grafana/custome-grafana-colmena
image_tag := 8.0

IMAGE := $(registry_host)/$(image_path):$(image_tag)

all: build push

# Build the Docker image with the proper tag
build:
	@echo "Building Docker image $(IMAGE)..."
	docker build \
	  --pull \
	  --tag $(IMAGE) \
	  .
	@echo "Build complete."

# Push requires prior login to the registry
push: login
	@echo "Pushing image $(IMAGE) to $(registry_host)..."
	docker push $(IMAGE)
	@echo "Push complete."

# Log in to Nexus (expects DOCKER_USER and DOCKER_PASS env vars)
login:
ifndef DOCKER_USER
	$(error DOCKER_USER is not set)
endif
ifndef DOCKER_PASS
	$(error DOCKER_PASS is not set)
endif
	@echo "Logging in to $(registry_host)..."
	echo "$(DOCKER_PASS)" | docker login $(registry_host) -u "$(DOCKER_USER)" --password-stdin
	@echo "Login successful."

# Optionally remove built image locally
clean:
	@echo "Cleaning up local image $(IMAGE)..."
	docker image rm -f $(IMAGE) || true
	@echo "Cleanup done."

