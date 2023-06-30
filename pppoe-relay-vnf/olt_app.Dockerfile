FROM golang:1.18.5 AS builder

WORKDIR /pppoe-vnf
COPY . .
RUN cd pppoe-relay-vnf/oltapp/mainapp && go build -o /oltapp

FROM ubuntu:latest
WORKDIR /
RUN apt update -y && apt install net-tools
COPY --from=builder /oltapp /
COPY ./pppoe-relay-vnf/olt_app.sh /