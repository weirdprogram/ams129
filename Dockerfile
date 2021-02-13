FROM golang:1.15.8-buster
RUN apt-get update -y && apt-get install -y python3 python3-pandas python3-redis redis-server