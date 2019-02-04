FROM golang

RUN go get -d -v github.com/lib/pq && \
    go get -d -v github.com/go-redis/redis && \
    go get -d -v github.com/tarantool/go-tarantool && \
    go get -d -v github.com/gorilla/mux

WORKDIR /go/src/
COPY testtask ./testtask

RUN go build testtask/cmd/tarantool

CMD ["go", "run", "testtask/cmd/tarantool"]
