FROM alpine:3.17.1

ADD ./flannel-network-health /flannel-network-health

ENTRYPOINT ["/flannel-network-health"]
CMD ["daemon"]
