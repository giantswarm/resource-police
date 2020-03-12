FROM quay.io/giantswarm/alpine:3.9-giantswarm

ADD ./resource-police /resource-police

ENTRYPOINT ["/resource-police"]
