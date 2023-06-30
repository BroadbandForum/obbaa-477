FROM ubuntu:latest

RUN apt -y update && apt install -y pppoe

WORKDIR /server

COPY pppoe-server-options /etc/ppp/pppoe-server-options
COPY chap-secrets /etc/ppp/chap-secrets

COPY server.sh .

# CMD ./server.sh