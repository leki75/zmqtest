FROM golang:1.18

RUN apt-get update && \
    apt-get install -y libczmq-dev && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/zmqtest

COPY . .

RUN go mod tidy

CMD ["go", "test", "-bench=.", "-benchmem", "./..."]
