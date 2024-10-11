FROM golang:1.18

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /hospital-service ./cmd

EXPOSE 8081

CMD ["./hospital-service"]
