#!/bin/bash

cd ./server

bash c.sh &

cd ../app

echo "installing modules..."
npm install

echo "running the app..."
bash c.sh &

echo "terminals are ready"
