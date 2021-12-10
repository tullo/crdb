#!/bin/bash
echo "Terminating $1 process by sending SIGTERM";
PID=$(ps u | grep "store=$1" | awk '{print $2}' | head -n1)
if kill ${PID}
then
  echo "Terminated ${PID}"
fi
