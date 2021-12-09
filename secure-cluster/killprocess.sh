#!/bin/bash
echo "Killing $1 process  by sending SIGKILL";
PID=$(ps u | grep "store=$1" | awk '{print $2}' | head -1)
if kill -9 ${PID}
then
  echo "Killed ${PID}"
fi
