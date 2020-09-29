FROM golang:1.13-alpine as builder

# Ultimate .dockerignore should protect
# from copying unnecessary data (.git/ dir, etc.)
COPY . .
RUN go install ./...

FROM alpine:3.10
RUN apk add --no-cache tzdata
CMD ["/usr/local/bin/wireguard-for-homies"]
COPY --from=builder /go/bin/wireguard-for-homies /usr/local/bin/wireguard-for-homies

ARG BUILD_DATE
ARG VCS_REF
ARG VERSION

LABEL org.label-schema.build-date $BUILD_DATE
LABEL org.label-schema.name "wireguard-for-homies"
LABEL org.label-schema.vcs-url "https://github.com/cordalace/wireguard-for-homies"
LABEL org.label-schema.vcs-ref $VCS_REF
LABEL org.label-schema.vendor "Azat Kurbanov <cordalace@gmail.com>"
LABEL org.label-schema.version $VERSION
LABEL org.label-schema.schema-version "1.0"

ENV VERSION $VERSION
ENV ENVIRONMENT production
