#!/bin/bash

rm -Rf /c-deployer/bin/canal.pid

sh /c-deployer/bin/startup.sh

while true;do echo beatheart;sleep 5;done;
