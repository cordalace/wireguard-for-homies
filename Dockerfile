FROM golang:1.15-buster as builder

WORKDIR /build

RUN apt-get update && apt-get install -y --no-install-recommends gcc libc-dev

# Ultimate .dockerignore should protect
# from copying unnecessary data (.git/ dir, etc.)
COPY . .
RUN go install -v ./...

FROM gcr.io/distroless/base-debian10
CMD ["/go/bin/wireguard-for-homies"]
COPY --from=builder /go/bin/wireguard-for-homies /go/bin/wireguard-for-homies

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
