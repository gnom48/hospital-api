FROM golang:1.23-alpine

WORKDIR /account

COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./config ./config

RUN go build -o ./account-service ./cmd/account

EXPOSE 8081

CMD ["sh", "-c", "sleep 5; while [ ! -f ./account-service ]; do sleep 1; done; ./account-service"]