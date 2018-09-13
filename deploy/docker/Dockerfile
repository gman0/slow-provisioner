FROM alpine:3.7

RUN apk add --no-cache ca-certificates

ADD slow-provisioner /slow-provisioner

ENTRYPOINT ["/slow-provisioner"]
