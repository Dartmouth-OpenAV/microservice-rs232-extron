FROM golang:latest

COPY source /go/src

ENV GOPATH=

WORKDIR /go/src/microservice-framework
RUN go mod init github.com/Dartmouth-OpenAV/microservice-framework
RUN go mod tidy

WORKDIR /go
# Change the module path for each microservice
RUN go mod init github.com/Dartmouth-OpenAV/microservice-rs232-extron
RUN go mod edit -replace github.com/Dartmouth-OpenAV/microservice-framework=./src/microservice-framework
RUN go mod tidy

WORKDIR /go/src
RUN go get -u
RUN go build -o /go/bin/microservice

# Use this entrypoint for a a fully functional docker image
ENTRYPOINT /go/bin/microservice
