FROM golang:1.23.0 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY validator.go ./
COPY ./router ./router
COPY ./store ./store
COPY ./config ./config
COPY ./email ./email
COPY ./jwt ./jwt
COPY ./middleware ./middleware
COPY ./util ./util
COPY ./views ./views
COPY ./tag ./tag

ARG VERSION
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o konbini -ldflags "-X github.com/juancwu/konbini/config.Version=${VERSION}"

FROM alpine AS runner
WORKDIR /go
COPY --from=builder /app/konbini ./

EXPOSE 3000

ENTRYPOINT ["./konbini"]
