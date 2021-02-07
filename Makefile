PROJECT_PKGS := $$(go list -mod vendor ./pkg/... | grep -v /vendor/)
APP_NAME ?= tc-cli

format: ## Format sources
	go fmt $(PROJECT_PKGS)

lint: ## Run linter on sources
	golangci-lint run --modules-download-mode=vendor ./pkg...

test: ## Run unit tests without coverage report
	go test -mod vendor -v ./...

build: ## Build application
	mkdir -p build
	go build -mod=vendor -o build/${APP_NAME} cmd/main.go

build-image:
	docker-compose build tc-cli
