# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

ADD . /go/src/github.com/mkoptik/qcmd-server-go
RUN cd /go/src/github.com/mkoptik/qcmd-server-go

# Get all dependecies
RUN go get -v ./...

# Install to /go/bin/qcmd-server-go
RUN go install github.com/mkoptik/qcmd-server-go

ENTRYPOINT /go/bin/qcmd-server-go

EXPOSE 8888