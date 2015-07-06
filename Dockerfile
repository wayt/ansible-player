# Pull base image.
FROM google/golang

# Install HG for go get
#RUN apt-get update && \
#    apt-get install -y mercurial curl git


ADD . /gopath/src/github.com/wayt/ansible-player
WORKDIR /gopath/src/github.com/wayt/ansible-player

RUN go get
RUN go install

ENTRYPOINT ["/gopath/bin/ansible-player"]

EXPOSE 8080
