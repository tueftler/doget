FROM debian:jessie

MAINTAINER Tueftler Group 

COPY ./doget ./.doget.yml /usr/bin/

VOLUME ["/app"]

WORKDIR /app

ENTRYPOINT ["/usr/bin/doget"]