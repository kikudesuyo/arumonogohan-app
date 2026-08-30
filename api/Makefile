SHELL := /bin/bash

PROJECT_ID ?= arumonogohan-app
REGION ?= asia-northeast1
SERVICE ?= arumonogohan-api
REPOSITORY ?= arumonogohan
RUNTIME_SERVICE_ACCOUNT ?= arumonogohan-runtime@$(PROJECT_ID).iam.gserviceaccount.com
COMMIT_SHA ?= $(shell git rev-parse --short=12 HEAD 2>/dev/null || echo latest)

GEMINI_PROJECT_ID ?= $(PROJECT_ID)
GEMINI_API_KEY_SECRET ?= GEMINI_API_KEY
LINE_CHANNEL_SECRET_SECRET ?= LINE_BOT_CHANNEL_SECRET
LINE_CHANNEL_TOKEN_SECRET ?= LINE_BOT_CHANNEL_TOKEN

.PHONY: format format-check lint test build docker-build ci image deploy

format:
	gofmt -w $$(find . -type f -name '*.go')

format-check:
	test -z "$$(gofmt -l .)"

lint:
	go vet ./...

test:
	go test ./...

build:
	go build -o /tmp/arumonogohan-api .

docker-build:
	docker build --tag arumonogohan-api:local .

ci: format-check lint test build docker-build

# deploy is intended to run from the CI workflow after Workload Identity Federation auth.
deploy:
	gcloud builds submit \
		--project="$(PROJECT_ID)" \
		--config=_cloudbuild.all.yaml \
		--substitutions=_COMMIT_SHA="$(COMMIT_SHA)",_REGION="$(REGION)",_SERVICE="$(SERVICE)",_RUNTIME_SERVICE_ACCOUNT="$(RUNTIME_SERVICE_ACCOUNT)",_GEMINI_PROJECT_ID="$(GEMINI_PROJECT_ID)",_GEMINI_API_KEY_SECRET="$(GEMINI_API_KEY_SECRET)",_LINE_CHANNEL_SECRET_SECRET="$(LINE_CHANNEL_SECRET_SECRET)",_LINE_CHANNEL_TOKEN_SECRET="$(LINE_CHANNEL_TOKEN_SECRET)" \
		.
