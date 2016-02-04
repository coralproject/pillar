# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.5.3

# Copy the local package files to the container's workspace.
COPY . /go/src/github.com/coralproject/pillar

# Go get all necessary packages
RUN go get github.com/gorilla/mux
RUN go get github.com/gorilla/handlers
RUN go get gopkg.in/mgo.v2
RUN go get gopkg.in/mgo.v2/bson
RUN go get gopkg.in/bluesuncorp/validator.v6

# Build & Install
RUN cd /go/src && go install github.com/coralproject/pillar/app/pillar

# Run the app
ENTRYPOINT /go/bin/pillar

# Document that the service listens on port 8080.
EXPOSE 8080

