FROM ubuntu:focal as dqlite-lib-builder
ARG TARGETPLATFORM
ARG TARGETARCH
ARG TARGETOS
ARG GO_VERSION
ARG DQLITE_VERSION=v1.9.0
ARG RAFT_VERSION=v0.11.2
ARG DEBIAN_FRONTEND="noninteractive"
ENV TZ=America/New_York
ENV LD_LIBRARY_PATH=/usr/local/lib
ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV DQLITE_VERSION=${DQLITE_VERSION}
ENV RAFT_VERSION=${RAFT_VERSION}
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH
ENV GO_VERSION=${GO_VERSION} OS=linux ARCH=${TARGETARCH}

RUN apt update -y && \
    apt install -y software-properties-common apt-utils git build-essential dh-autoreconf pkg-config libuv1-dev libsqlite3-dev liblz4-1 liblz4-dev

WORKDIR /opt/raft

RUN git clone --depth 1 -b $RAFT_VERSION -v  https://github.com/canonical/raft.git ./ && \
    autoreconf -i && ./configure && make && make install

WORKDIR /opt/dqlite

RUN git clone --depth 1 -b $DQLITE_VERSION https://github.com/canonical/dqlite.git ./ && \
    autoreconf -i && ./configure && make && make install

WORKDIR /opt/golang

RUN curl -o go.tar.gz https://dl.google.com/go/go$GO_VERSION.$OS-$TARGETARCH.tar.gz && \
    rm -rf /usr/local/go && \
    tar -C /usr/local -xzf go.tar.gz