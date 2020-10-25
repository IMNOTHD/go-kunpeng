#!/bin/bash

rm -Rf /c-deployer/bin/canal.pid

sh /c-deployer/bin/startup.sh

while true;do echo "`date '+%Y-%m-%d %H:%M:%S'`: heartbeat";sleep 5;done;
