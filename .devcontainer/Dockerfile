FROM ubuntu:24.04

RUN apt-get update && apt-get install -yy -q wget
RUN wget https://github.com/caddyserver/xcaddy/releases/download/v0.4.2/xcaddy_0.4.2_linux_amd64.deb && \
    dpkg -i xcaddy_0.4.2_linux_amd64.deb && \
    rm xcaddy_0.4.2_linux_amd64.deb