# Copyright 2020 Jayson Vibandor. All rights reserved

# STEP 1: Build the executable
FROM golang:1.13.7-alpine3.11 as builder
RUN apk update && apk upgrade && \
    apk --update add git gcc make
WORKDIR /go/src/gophr
COPY . .
ENV GO111MODULE on
RUN make build-api -s

# STEP 2: Distribute the executable
FROM alpine:3.7
RUN apk update && \
    apk upgrade && \
    apk add --no-cache bash && \
    apk add --no-cache ca-certificates tzdata && \
    apk add --no-cache openssh
RUN set -ex && apk add --no-cache --virtual bash musl-dev openssl
WORKDIR /home/gophr/
ENV HOME /home/
EXPOSE 8080
COPY --from=builder /go/src/gophr/bin/gophr.engine /home/gophr/
RUN chmod +x /home/gophr/gophr.engine
CMD /home/gophr/gophr.engine


