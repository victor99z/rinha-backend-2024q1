#!/bin/bash

echo "Build and run dev env"

docker container rm -f $(docker container ps -aq)
docker-compose -f docker-compose-dev.yml up --build
