FROM golang:1.25.3-alpine AS builder
WORKDIR /src
COPY . .
RUN apk add --no-cache musl-dev build-base git
RUN go mod download
RUN CGO_ENABLED=1 CC=gcc GOOS=linux go build -o agent -a -ldflags '-linkmode external -extldflags "-static"' cmd/agent/main.go
 
FROM scratch
COPY --from=builder /src/agent ./agent
EXPOSE 8080
CMD ["./agent"]
