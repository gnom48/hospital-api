FROM golang:1.23-alpine

WORKDIR /hospital

COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./config ./config

RUN go build -o ./hospital-service ./cmd/hospital

EXPOSE 8081

CMD ["sh", "-c", "sleep 5; while [ ! -f ./hospital-service ]; do sleep 1; done; ./hospital-service"]