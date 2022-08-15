FROM alpine:3.16.2

ADD ./flannel-network-health /flannel-network-health

ENTRYPOINT ["/flannel-network-health"]
CMD ["daemon"]
