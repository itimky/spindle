FROM golang:1.21.6-alpine3.19 AS builder
WORKDIR /src
RUN go env -w GOMODCACHE=/root/.cache/go-build
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build go mod download
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build go build -o /go/bin/spindle ./cmd/spindle

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /go/bin/spindle /app
CMD ["/app/spindle"]
