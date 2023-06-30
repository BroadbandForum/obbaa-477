#!/bin/bash
ifconfig eth0 0
pppd call provider
tail -f /dev/null