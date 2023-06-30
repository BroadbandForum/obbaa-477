#!/bin/bash
pppoe-server persist debug -r -S server -C server -T 10 -I eth0 -L $1 -R $2
tail -f /dev/null