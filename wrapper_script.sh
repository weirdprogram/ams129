#!/bin/bash

# https://docs.docker.com/config/containers/multi-service_container/
# Start the first process
/usr/bin/redis-server --daemonize yes
status=$?
if [ $status -ne 0 ]; then
  echo "Failed to start /usr/bin/redis-server: $status"
  exit $status
fi

# # Start the second process
# ./my_second_process -D
# status=$?
# if [ $status -ne 0 ]; then
#   echo "Failed to start my_second_process: $status"
#   exit $status
# fi

# Naive check runs checks once a minute to see if either of the processes exited.
# This illustrates part of the heavy lifting you need to do if you want to run
# more than one service in a container. The container exits with an error
# if it detects that either of the processes has exited.
# Otherwise it loops forever, waking up every 60 seconds

while sleep 15; do
  ps aux |grep redis-server |grep -q -v grep
  PROCESS_1_STATUS=$?
  # ps aux |grep my_second_process |grep -q -v grep
  # PROCESS_2_STATUS=$?
  # If the greps above find anything, they exit with 0 status
  # If they are not both 0, then something is wrong
  # if [ $PROCESS_1_STATUS -ne 0 -o $PROCESS_2_STATUS -ne 0 ]; then
  if [ $PROCESS_1_STATUS -ne 0 ]; then
    echo "One of the processes has already exited."
    exit 1
  fi
done