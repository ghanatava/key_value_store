FROM golang:1.18 as builder
WORKDIR /src/kvs
COPY . /src/kvs
RUN GOOS=linux go build -o kvs

FROM ubuntu
COPY --from=builder /src/kvs /src/kvs
WORKDIR /src/kvs

CMD ["./kvs"] 