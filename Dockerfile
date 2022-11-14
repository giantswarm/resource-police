FROM quay.io/giantswarm/alpine:3.16.3-giantswarm

ADD ./resource-police /resource-police

ENTRYPOINT ["/resource-police"]
