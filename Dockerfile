FROM alpine:3.6

ADD ./flannel-network-health /flannel-network-health

ENTRYPOINT ["/main"]
