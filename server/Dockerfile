# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

ENV MONGODB_URL mongodb://localhost:27017/coral

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/coralproject/pillar

# Go get all necessary packages
RUN go get github.com/gorilla/mux
RUN go get gopkg.in/mgo.v2
RUN go get gopkg.in/mgo.v2/bson

# Build & Install
RUN cd /go/src && go install github.com/coralproject/pillar

# Run the app
ENTRYPOINT /go/bin/pillar

# Document that the service listens on port 8080.
EXPOSE 8080
