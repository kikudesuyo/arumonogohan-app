SHELL := /bin/bash

PROJECT_ID ?= arumonogohan-app
REGION ?= asia-northeast1
SERVICE ?= arumonogohan-api
REPOSITORY ?= arumonogohan
IMAGE_TAG ?= $(shell git rev-parse --short=12 HEAD 2>/dev/null || echo latest)
IMAGE := $(REGION)-docker.pkg.dev/$(PROJECT_ID)/$(REPOSITORY)/$(SERVICE):$(IMAGE_TAG)
RUNTIME_SERVICE_ACCOUNT ?= arumonogohan-runtime@$(PROJECT_ID).iam.gserviceaccount.com

GEMINI_PROJECT_ID ?= $(PROJECT_ID)
GEMINI_API_KEY_SECRET ?= gemini-api-key
LINE_CHANNEL_SECRET_SECRET ?= line-bot-channel-secret
LINE_CHANNEL_TOKEN_SECRET ?= line-bot-channel-token

.PHONY: format format-check lint test build docker-build ci image deploy

format:
	gofmt -w $$(find api -type f -name '*.go')

format-check:
	test -z "$$(gofmt -l api)"

lint:
	go vet ./...

test:
	go test ./...

build:
	go build -o /tmp/arumonogohan-api ./api

docker-build:
	docker build --tag arumonogohan-api:local .

ci: format-check lint test build docker-build

# deploy is intended to run from the CI workflow after Workload Identity Federation auth.
image:
	gcloud builds submit \
		--project="$(PROJECT_ID)" \
		--tag="$(IMAGE)" \
		.

deploy: image
	gcloud run deploy "$(SERVICE)" \
		--image="$(IMAGE)" \
		--region="$(REGION)" \
		--project="$(PROJECT_ID)" \
		--platform="managed" \
		--service-account="$(RUNTIME_SERVICE_ACCOUNT)" \
		--allow-unauthenticated \
		--port=8080 \
		--cpu=1 \
		--memory=256Mi \
		--min-instances=0 \
		--max-instances=1 \
		--update-env-vars="GEMINI_PROJECT_ID=$(GEMINI_PROJECT_ID)" \
		--set-secrets="GEMINI_API_KEY=$(GEMINI_API_KEY_SECRET):latest,LINE_BOT_CHANNEL_SECRET=$(LINE_CHANNEL_SECRET_SECRET):latest,LINE_BOT_CHANNEL_TOKEN=$(LINE_CHANNEL_TOKEN_SECRET):latest"
