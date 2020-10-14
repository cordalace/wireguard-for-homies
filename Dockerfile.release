FROM golang:1.15-buster as builder

WORKDIR /build

RUN apt-get update && apt-get install -y --no-install-recommends gcc libc-dev

# Ultimate .dockerignore should protect
# from copying unnecessary data (.git/ dir, etc.)
COPY . .
RUN go install ./...

FROM gcr.io/distroless/base-debian10
CMD ["/go/bin/wireguard-for-homies"]
COPY --from=builder /go/bin/wireguard-for-homies /go/bin/wireguard-for-homies
