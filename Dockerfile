FROM alpine:3.14.0

ADD ./flannel-network-health /flannel-network-health

ENTRYPOINT ["/flannel-network-health"]
CMD ["daemon"]
