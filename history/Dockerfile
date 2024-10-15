FROM golang:1.23-alpine

WORKDIR /history

COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./config ./config

RUN go build -o ./history-service ./cmd/history

EXPOSE 8081

CMD ["sh", "-c", "sleep 5; while [ ! -f ./history-service ]; do sleep 1; done; ./history-service"]