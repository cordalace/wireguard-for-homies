FROM gcr.io/distroless/base-debian10
CMD ["/go/bin/wireguard-for-homies"]
COPY wireguard-for-homies /go/bin/wireguard-for-homies
