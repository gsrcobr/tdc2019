#!/bin/bash

# Start server
pgrep -f "bind 127.0.0.1" > ./server.pid
if [ $? != 0 ]; then
  consul agent -dev -bind 127.0.0.1 -config-dir=./config >> ./server.log 2>&1 &
  echo $! > ./server.pid
  echo "Client 1 / Server started! PID is $(cat server.pid)"
else
  echo "Client 1 / Server already running! PID is $(cat server.pid)"
fi

# Start nodes
ip -4 addr show |grep lo0 |fgrep -q 127.0.0.2
if [ $? != 0 ]; then
  sudo ip addr add 127.0.0.2/8 dev lo0
  echo "Interface address for client 2 added!"
fi
pgrep -f "bind 127.0.0.2" > ./client2.pid
if [ $? != 0 ]; then
  consul agent -join 127.0.0.1 -bind 127.0.0.2 -client 127.0.0.2 -config-dir=./config2 >> ./client2.log 2>&1 &
  echo $! > ./client2.pid
  echo "Client 2 started! PID is $(cat client2.pid)"
else
  echo "Client 2 already running! PID is $(cat client2.pid)"
fi

ip -4 addr show |grep lo0 |fgrep -q 127.0.0.3
if [ $? != 0 ]; then
  sudo ip addr add 127.0.0.3/8 dev lo0
  echo "Interface address for client 3 added!"
fi
pgrep -f "bind 127.0.0.3" > ./client3.pid
if [ $? != 0 ]; then
  consul agent -join 127.0.0.1 -bind 127.0.0.3 -client 127.0.0.3 -config-dir=./config3 >> ./client3.log 2>&1 &
  echo $! > ./client3.pid
  echo "Client 3 started! PID is $(cat client3.pid)"
else
  echo "Client 3 already running! PID is $(cat client3.pid)"
fi

exit 0
