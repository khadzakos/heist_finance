DOCKER_COMPOSE = docker-compose
DOCKER = docker
EXCHANGES = binance bybit coinbase okx 
CONNECTOR_IMAGES = $(addsuffix -connector,$(EXCHANGES))  # Generates: bybit-connector coinbase-connector okx-connector
PREPROCESSOR_IMAGES = $(addsuffix -preprocessor,$(EXCHANGES))  # Generates: bybit-preprocessor coinbase-preprocessor okx-preprocessor
ALL_IMAGES = $(CONNECTOR_IMAGES) $(PREPROCESSOR_IMAGES)
PROJECT_NAME = heist
IMAGE_PREFIX = $(PROJECT_NAME)/

.PHONY: all
all: build up

.PHONY: build
build: $(ALL_IMAGES)

.PHONY: $(CONNECTOR_IMAGES)
$(CONNECTOR_IMAGES):
	$(eval EXCHANGE := $(subst -connector,,$@))  # Extract exchange name (e.g., bybit from bybit-connector)
	@if [ -z "$$($(DOCKER) images -q $(IMAGE_PREFIX)$@)" ]; then \
		echo "Image $(IMAGE_PREFIX)$@ not found, building..."; \
		$(DOCKER) build -t $(IMAGE_PREFIX)$@ --build-arg EXCHANGE=$(EXCHANGE) -f connector/Dockerfile connector; \
	else \
		echo "Image $(IMAGE_PREFIX)$@ already exists, skipping build."; \
	fi

.PHONY: $(PREPROCESSOR_IMAGES)
$(PREPROCESSOR_IMAGES):
	$(eval EXCHANGE := $(subst -preprocessor,,$@))  # Extract exchange name (e.g., bybit from bybit-preprocessor)
	@if [ -z "$$($(DOCKER) images -q $(IMAGE_PREFIX)$@)" ]; then \
		echo "Image $(IMAGE_PREFIX)$@ not found, building..."; \
		$(DOCKER) build -t $(IMAGE_PREFIX)$@ --build-arg EXCHANGE=$(EXCHANGE) -f preprocessor/Dockerfile preprocessor; \
	else \
		echo "Image $(IMAGE_PREFIX)$@ already exists, skipping build."; \
	fi

.PHONY: up
up:
	$(DOCKER_COMPOSE) up -d

.PHONY: down
down:
	$(DOCKER_COMPOSE) down

.PHONY: clean
clean:
	$(DOCKER_COMPOSE) down --rmi all
	@for img in $(ALL_IMAGES); do \
		if [ -n "$$($(DOCKER) images -q $(IMAGE_PREFIX)$$img)" ]; then \
			$(DOCKER) rmi $(IMAGE_PREFIX)$$img; \
		fi; \
	done

.PHONY: rebuild
rebuild: clean build up

.PHONY: images
images:
	@echo "Docker images to be built:"
	@for img in $(ALL_IMAGES); do \
		echo "  - $(IMAGE_PREFIX)$$img"; \
	done