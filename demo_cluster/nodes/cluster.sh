#!/usr/bin/env bash
ps -uax|grep nodes|awk '{print $2}'|xargs kill -9 

./nodes 1  --path=../../config/demo-cluster.json --node=0.0.1.1 >> master.log 2>&1 &
./nodes 2  --path=../../config/demo-cluster.json --node=0.0.2.1 >> center.log 2>&1 &
./nodes 3  --path=../../config/demo-cluster.json --node=0.0.3.1 >> web.log 2>&1 &
./nodes 4  --path=../../config/demo-cluster.json --node=0.0.4.1 >> gate.log 2>&1 &
./nodes 5  --path=../../config/demo-cluster.json --node=1.1.5.1 >> game.log 2>&1 &

