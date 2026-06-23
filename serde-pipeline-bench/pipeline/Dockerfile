# Builds the producer and consumer into a tiny runtime image. Build context is
# the repo root (see docker-compose.yml) so the whole module is available.
FROM golang:1.25 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/producer ./pipeline/cmd/producer
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/consumer ./pipeline/cmd/consumer

FROM gcr.io/distroless/static-debian12
COPY --from=build /out/producer /producer
COPY --from=build /out/consumer /consumer
# Default entrypoint is the consumer; producer is invoked explicitly via compose.
ENTRYPOINT ["/consumer"]
