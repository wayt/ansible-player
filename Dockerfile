# Pull base image.
FROM ubuntu:trusty

# Install ansible stuff
RUN apt-get install -y software-properties-common
RUN apt-add-repository ppa:ansible/ansible
RUN apt-get update
RUN apt-get install -y ansible

RUN apt-get update -y && apt-get install --no-install-recommends -y -q curl build-essential ca-certificates git mercurial bzr
RUN mkdir /goroot && curl https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz | tar xvzf - -C /goroot --strip-components=1
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
