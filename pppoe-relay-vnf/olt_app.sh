#!/bin/bash

intf_of() {
    ifconfig | grep -B1 $1 | grep -o "^\w*"
}

if [[ "$RUN_IMMEDIATELY" = "true" ]]
then
    client_if=$(intf_of $CLIENT_ADDR)
    server_if=$(intf_of $SERVER_ADDR)
    ./oltapp -client_if $client_if -server_if $server_if -vnf-addr:port $VNF_ADDR:$VNF_PORT
else 
    tail -f /dev/null
fi

