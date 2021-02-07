FROM golang:1.15.5

LABEL description="Teamcity CLI docker image" \
      maintainer="Ilya Golovchenko (furiousassault)" \
      source="https://github.com/furiousassault/tc-cli"

WORKDIR /app
COPY . /app

RUN make build

ENTRYPOINT ["build/tc-cli"]