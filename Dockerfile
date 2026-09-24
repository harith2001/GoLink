FROM golang:1.26-bookworm AS build
WORKDIR /src
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 go build -o /server ./cmd/server

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=build /server /server
COPY backend/migrations ./migrations
COPY seed ./seed
ENV TRANSPORT=http
ENV SEED_DIR=/app/seed
ENV MIGRATIONS_DIR=/app/migrations
EXPOSE 8080
CMD ["/server"]
