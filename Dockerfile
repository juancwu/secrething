FROM golang:1.22.0 as builder
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

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o konbini

FROM alpine as runner
WORKDIR /go
COPY --from=builder /app/konbini ./

EXPOSE 3000

ENV APP_ENV="production"

ENTRYPOINT ["./konbini"]
