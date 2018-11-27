# Pull base image.
FROM ubuntu:bionic

# Install ansible stuff
RUN apt-get update && \
    apt-get install -y software-properties-common && \
    apt-add-repository ppa:ansible/ansible && \
    apt-get update && \
    apt-get install -y -q ansible

RUN apt-get install --no-install-recommends -y -q \
        curl \
        build-essential \
        ca-certificates \
        git \
        mercurial \
        bzr \
        && rm -rf /var/lib/apt/lists/*
RUN mkdir /goroot && curl https://storage.googleapis.com/golang/go1.11.2.linux-amd64.tar.gz | tar xvzf - -C /goroot --strip-components=1
RUN mkdir /gopath

ENV GOROOT /goroot
ENV GOPATH /gopath
ENV PATH $PATH:$GOROOT/bin:$GOPATH/bin

ADD . /gopath/src/github.com/wayt/ansible-player
WORKDIR /gopath/src/github.com/wayt/ansible-player

RUN go get
RUN go install

ENTRYPOINT ["/gopath/bin/ansible-player"]

EXPOSE 8080
