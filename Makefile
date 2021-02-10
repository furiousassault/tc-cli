PROJECT_PKGS := $$(go list -mod vendor ./pkg/... | grep -v /vendor/)
APP_NAME ?= tc-cli

format:
	go fmt $(PROJECT_PKGS)

lint:
	golangci-lint run --modules-download-mode=vendor ./pkg...

test:
	go test -mod vendor -cover -v ./...

build:
	mkdir -p build
	go build -mod=vendor -o build/${APP_NAME} cmd/main.go

build-image:
	docker-compose build tc-cli
