#--- Build stage
FROM golang:1.18.1-stretch AS go-builder

WORKDIR /src

COPY . /src/

RUN make build CGO_ENABLED=0

#--- Image stage
FROM alpine:3.15

COPY --from=go-builder /src/target/dist/cosmos-faucet /usr/bin/cosmos-faucet

WORKDIR /opt

ENTRYPOINT ["/usr/bin/cosmos-faucet"]
