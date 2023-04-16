#!/bin/sh
docker build -t forum -f Dockerfile .
docker container run --rm -it -p 5050:8080 forum