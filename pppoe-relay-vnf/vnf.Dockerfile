FROM golang:1.18.5 AS builder

WORKDIR /pppoe-vnf
COPY . .
RUN cd pppoe-relay-vnf/vnf/mainapp && go build -o /vnf

FROM ubuntu:20.04
WORKDIR /

COPY --from=builder /vnf /

CMD ./vnf

