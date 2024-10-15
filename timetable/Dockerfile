FROM golang:1.23-alpine

WORKDIR /timetable

COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./config ./config

RUN go build -o ./timetable-service ./cmd/timetable

EXPOSE 8081

CMD ["sh", "-c", "sleep 5; while [ ! -f ./timetable-service ]; do sleep 1; done; ./timetable-service"]