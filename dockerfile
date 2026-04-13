FROM ubuntu:18.04

LABEL org.opencontainers.image.source=https://github.com/karimmdjdb/mqtt-mapper

RUN mkdir -p kubeedge

COPY  /build/main kubeedge/
COPY ./config.yaml kubeedge/

WORKDIR kubeedge
