module httpserver

go 1.16

require (
	et-protobuf3 v0.0.0
	github.com/gin-gonic/gin v1.8.1
	github.com/go-micro/plugins/v4/registry/etcd v1.1.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	// google.golang.org/grpc v1.48.0
	go-micro.dev/v4 v4.8.0
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	et-config v0.0.0
)

replace et-protobuf3 => ../et-protobuf3
replace et-config => ../et-config