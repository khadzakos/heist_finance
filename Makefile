DOCKER_COMPOSE = docker-compose
DOCKER = docker
EXCHANGES = binance bybit coinbase okx moex nyse nasdaq lseg
CONNECTOR_IMAGES = $(addsuffix -connector,$(EXCHANGES))
PREPROCESSOR_IMAGES = $(addsuffix -preprocessor,$(EXCHANGES)) 
ALL_IMAGES = $(CONNECTOR_IMAGES) $(PREPROCESSOR_IMAGES) 
PROJECT_NAME = heist
IMAGE_PREFIX = $(PROJECT_NAME)/

.PHONY: all
all: build up ui

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

.PHONY: ui
ui:
	$(DOCKER) build -t $(IMAGE_PREFIX)ui -f ui/Dockerfile ui
	$(DOCKER) run -d -p 3000:3000 --name ui $(IMAGE_PREFIX)ui

.PHONY: up
up:
	$(DOCKER_COMPOSE) up --build -d 

.PHONY: down
down:
	$(DOCKER_COMPOSE) down

.PHONY: rebuild
rebuild: build up

.PHONY: images
images:
	@echo "Docker images to be built:"
	@for img in $(ALL_IMAGES); do \
		echo "  - $(IMAGE_PREFIX)$$img"; \
	done

.PHONY: test test-all test-cleaner test-preprocessor test-controller test-request-service

test: test-all

test-all: test-cleaner test-preprocessor test-controller test-request-service

test-cleaner:
	cd cleaner && make test

test-preprocessor:
	cd preprocessor && make test

test-controller:
	cd controller && make test

test-request-service:
	cd request-service && make test