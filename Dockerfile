FROM alpine:3.13.5

ADD ./flannel-network-health /flannel-network-health

ENTRYPOINT ["/flannel-network-health"]
CMD ["daemon"]
