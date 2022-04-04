FROM alpine:3.15.3

ADD ./flannel-network-health /flannel-network-health

ENTRYPOINT ["/flannel-network-health"]
CMD ["daemon"]
