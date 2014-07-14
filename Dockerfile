FROM ubuntu:14.04
MAINTAINER Devashish <devashish.86@gmail.com>

RUN mkdir -p /opt/nstop
ADD . /opt/nstop
RUN "cd /opt/nstop && make"
