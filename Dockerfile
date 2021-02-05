FROM golang:1.15.5

LABEL description="Teamcity CLI docker image" \
      maintainer="Ilya Golovchenko (furiousassault)" \
      source="https://github.com/furiousassault/tc-cli"

WORKDIR /app
COPY . /app

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -tags netgo --ldflags '-extldflags "-static"' -o /tc-cli cmd/main.go
RUN rm -rf /app

ENTRYPOINT ["/tc-cli"]
