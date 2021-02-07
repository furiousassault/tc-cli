PROJECT_PKGS := $$(go list ./pkg/... | grep -v /vendor/)
APP_NAME ?= tc-cli

IMAGE_VERSION ?= latest
IMAGE_PATH = "furiousassault/tc-client:latest"


format: ## Format sources
	go fmt $(PROJECT_PKGS)

lint: ## Run linter on sources
	GO111MODULE=on golangci-lint run --modules-download-mode=vendor ./pkg...

test: ## Run unit tests without coverage report
	GO111MODULE=on go test -mod vendor -v ./...

build: ## Build application
	mkdir -p build
	GO111MODULE=on go build -mod vendor -tags netgo --ldflags '-extldflags "-static"' -o build/${APP_NAME} main.go

