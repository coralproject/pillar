# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.5.2

ENV PILLAR_HOME /opt/pillar

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/coralproject/pillar/server

# Go get all necessary packages
RUN go get github.com/gorilla/mux
RUN go get gopkg.in/mgo.v2
RUN go get gopkg.in/mgo.v2/bson

# Build & Install
RUN cd /go/src && go install github.com/coralproject/pillar/server

# Run the app
ENTRYPOINT /go/bin/server

# Document that the service listens on port 8080.
EXPOSE 8080
