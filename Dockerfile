FROM gsoci.azurecr.io/giantswarm/alpine:3.20.3-giantswarm

ARG TARGETARCH
COPY resource-police-linux-${TARGETARCH} /resource-police

ENTRYPOINT ["/resource-police"]
