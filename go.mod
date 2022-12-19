module github.com/giantswarm/resource-police

go 1.16

require (
	github.com/giantswarm/microerror v0.4.0
	github.com/giantswarm/micrologger v1.0.0
	github.com/google/go-cmp v0.5.8
	github.com/kr/text v0.2.0 // indirect
	github.com/prometheus/client_golang v1.14.0
	github.com/prometheus/common v0.39.0
	github.com/spf13/cobra v1.3.0
)

replace (
	github.com/coreos/etcd v3.3.13+incompatible => github.com/etcd-io/etcd v3.3.25+incompatible
	github.com/dgrijalva/jwt-go => github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1
	github.com/gogo/protobuf v1.2.1 => github.com/gogo/protobuf v1.3.2
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.2
)
