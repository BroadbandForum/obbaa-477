FROM ubuntu:latest

RUN apt -y update && apt install -y pppoe net-tools
WORKDIR /client

COPY provider /etc/ppp/peers/provider

COPY client.sh .