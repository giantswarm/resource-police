FROM quay.io/giantswarm/alpine:3.16.2-giantswarm

ADD ./resource-police /resource-police

ENTRYPOINT ["/resource-police"]
