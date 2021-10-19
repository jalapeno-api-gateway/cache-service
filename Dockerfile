FROM golang:1.15-alpine AS build_base

RUN apk add --no-cache git
WORKDIR /tmp/cache-service

COPY go.mod .
COPY go.sum .
RUN go clean --modcache
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o ./out/cache-service .

FROM scratch
LABEL maintainer="Julian Klaiber"

COPY --from=build_base /tmp/cache-service/out/cache-service /usr/bin/cache-service

ENTRYPOINT [ "/usr/bin/cache-service" ]


