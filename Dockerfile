FROM gsoci.azurecr.io/giantswarm/alpine:3.20.3-giantswarm

ADD ./resource-police /resource-police

ENTRYPOINT ["/resource-police"]
