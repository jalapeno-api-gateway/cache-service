# no alpien due to special compilation flags
FROM golang:1.18 AS build_base

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

# Skaffold passes in debug-oriented compiler flags
ARG SKAFFOLD_GO_GCFLAGS
RUN go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -o /app/cache-service .


FROM gcr.io/distroless/base

# Definition of this variable is used by 'skaffold debug' to identify a golang binary.
# Default behavior - a failure prints a stack trace for the current goroutine.
# See https://golang.org/pkg/runtime/
ENV GOTRACEBACK=single

WORKDIR /app
COPY --from=build_base /app/cache-service /app/cache-service
ENTRYPOINT ["/app/cache-service"]