module github.com/giantswarm/resource-police

go 1.13

require (
	github.com/giantswarm/gsclientgen/v2 v2.1.0
	github.com/giantswarm/microerror v0.3.0
	github.com/giantswarm/micrologger v0.5.0
	github.com/go-openapi/runtime v0.19.26
	github.com/go-openapi/strfmt v0.20.1
	github.com/google/go-cmp v0.5.5
	github.com/hako/durafmt v0.0.0-20200710122514-c0fb7b4da026
	github.com/spf13/cobra v1.1.3
	sigs.k8s.io/yaml v1.2.0
)

replace github.com/coreos/etcd v3.3.13+incompatible => github.com/coreos/etcd v3.3.25+incompatible
