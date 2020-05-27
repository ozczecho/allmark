FROM golang:latest
MAINTAINER Mike <ozczecho@yahoo.com>

# Install pandoc for RTF conversion
RUN apt-get update && apt-get install -qy pandoc

# Build
COPY . /go/src/allmark

RUN mkdir /cli

WORKDIR /go/src/allmark/cli
RUN go get ./
RUN env GOOS=linux GOARCH=amd64 go build -o /cli/allmark

# Data
RUN mkdir /data
ADD . /data

VOLUME ["/data"]

EXPOSE 33001

CMD ["/cli/allmark", "serve", "/data"]
