####MULTI-STAGE BUILD####
FROM golang:1.13.5-alpine3.10 as builder
RUN apk update && apk upgrade && \
    apk --update add git gcc make
WORKDIR /go/src/github.com/jayvib/gophr
COPY . .
ENV GO111MODULE on
RUN make build-client

# Distribute
FROM alpine:latest
RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    apk add --no-cache bash && \
    mkdir /gophr
WORKDIR /gophr
EXPOSE 8080
COPY --from=builder /go/src/github.com/jayvib/gophr/bin/gophr.client /gophr
COPY --from=builder /go/src/github.com/jayvib/gophr/gophr.client/templates /gophr/templates
COPY --from=builder /go/src/github.com/jayvib/gophr/gophr.client/assets /gophr/assets
CMD ["/gophr/gophr.client"]
