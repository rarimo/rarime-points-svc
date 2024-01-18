FROM golang:1.20-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/rarimo/points-svc
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/points-svc /go/src/github.com/rarimo/points-svc


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/points-svc /usr/local/bin/points-svc
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["points-svc"]
