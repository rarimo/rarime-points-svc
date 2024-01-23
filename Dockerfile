FROM golang:1.21-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/rarimo/rarime-points-svc
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/rarime-points-svc /go/src/github.com/rarimo/rarime-points-svc


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/rarime-points-svc /usr/local/bin/rarime-points-svc
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["rarime-points-svc"]
