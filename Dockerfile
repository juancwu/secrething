FROM golang:1.23.0 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY ./common ./common
COPY ./handler ./handler
COPY ./middleware ./middleware
COPY ./service ./service
COPY ./store ./store
COPY ./types ./types

ARG VERSION
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o konbini -ldflags "-X github.com/juancwu/konbini/config.Version=${VERSION}"

FROM alpine AS runner
WORKDIR /go
COPY --from=builder /app/konbini ./

EXPOSE 3000

ENTRYPOINT ["./konbini"]
