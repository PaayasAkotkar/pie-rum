#!/bin/bash

export LANG=en_US.UTF-8

echo "hold on right installing the pie-rum sdk 😃"
go mod tidy

echo "just a sec running the server 🤗"
go build -o app

echo "now running the file 🌟"
./app

echo "the pie-rum server started 🤩"