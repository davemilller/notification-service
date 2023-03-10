COMPOSE_NAME = notification-service

LOCAL_COMPOSE_FILE ?= deploy/local/docker-compose.yaml
LOCAL_PROJECT_NAME := $(COMPOSE_NAME)-local

.PHONY: build
build:
	docker-compose -p $(LOCAL_PROJECT_NAME) -f $(LOCAL_COMPOSE_FILE) build --progress=plain

.PHONY: up
up:
	POSTGRES_DATA=postgres-data docker-compose -p $(LOCAL_PROJECT_NAME) -f $(LOCAL_COMPOSE_FILE) up

.PHONY: down
down:
	docker-compose -p $(LOCAL_PROJECT_NAME) -f $(LOCAL_COMPOSE_FILE) down

.PHONY: local
local: build up

.PHONY: test
test:
	make -C ./app test

.PHONY: test-integration
test-integration:
	make -C ./app test-integration

.PHONY: lint-test
lint-test:
	make -C ./app lint
	make -C ./app test

YELLOW := "\e[1;33m"
GREEN := "\e[1;32m"
NC := "\e[0m"

INFO := @bash -c '\
	printf $(GREEN); \
	echo "=> $$1"; \
	printf $(NC)' VALUE