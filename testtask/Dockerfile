FROM golang

RUN go get -d -v github.com/lib/pq && \
    go get -d -v github.com/go-redis/redis && \
    go get -d -v github.com/tarantool/go-tarantool

WORKDIR /go/src/
COPY testtask ./testtask

RUN go build testtask/cmd/tarantool

CMD ["go", "run", "testtask/cmd/tarantool"]
