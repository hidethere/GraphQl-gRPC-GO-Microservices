
FROM golang:1.25

WORKDIR /go/src/github.com/hidethere/GraphQl-gRPC-GO-Microservices

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /usr/local/bin/app ./order/cmd/order

EXPOSE 8080
CMD ["app"]
