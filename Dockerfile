FROM gsoci.azurecr.io/giantswarm/alpine:3.18.3-giantswarm

ADD ./resource-police /resource-police

ENTRYPOINT ["/resource-police"]
