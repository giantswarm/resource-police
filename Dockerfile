FROM quay.io/giantswarm/alpine:3.17.1-giantswarm

ADD ./resource-police /resource-police

ENTRYPOINT ["/resource-police"]
