FROM golang:1.18-alpine AS buildenv
WORKDIR /src
ADD . /src

RUN GOOS=linux go build -o ./out/app-for-HR ./cmd/main.go

RUN chmod +x ./out/app-for-HR

FROM alpine:latest
WORKDIR /app
COPY --from=buildenv /src/out/app-for-HR .

#### Local application port
EXPOSE 9090

ENTRYPOINT ["/app/app-for-HR"]