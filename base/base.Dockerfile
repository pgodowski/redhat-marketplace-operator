FROM registry.access.redhat.com/ubi8/ubi:latest as golang
ARG TARGETPLATFORM
ARG TARGETARCH
ARG TARGETOS
ARG GO_VERSION
ARG ARCH=${TARGETARCH}
ARG OS=${TARGETOS:-linux}
ENV VERSION=${GO_VERSION} OS=${OS} ARCH=${ARCH}

RUN dnf -y install git make yum gzip \
  && dnf update \
  && yum -y update-minimal --security --sec-severity=Important --sec-severity=Critical \
  && yum clean all \
  && dnf clean all \
  && rm -rf /var/cache/yum


RUN curl -o go.tar.gz https://dl.google.com/go/go$VERSION.$OS-$ARCH.tar.gz && \
  rm -rf /usr/local/go && \
  sha256sum go.tar.gz && \
  tar -C /usr/local -xzf go.tar.gz && \
  echo 'PATH=$PATH:/usr/local/go/bin' >> /etc/profile && \
  echo 'PATH=$PATH:/usr/local/go/bin' >> $HOME/.profile
