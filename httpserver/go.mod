module httpserver

go 1.16

require (
	et-config v0.0.0
	et-protobuf3 v0.0.0
	github.com/gin-gonic/gin v1.8.1
	github.com/go-micro/plugins/v4/registry/etcd v1.1.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/go-cmp v0.5.8 // indirect
	// google.golang.org/grpc v1.48.0
	go-micro.dev/v4 v4.8.0
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/zap v1.21.0
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/net v0.0.0-20211015210444-4f30a5c0130f // indirect
	golang.org/x/sys v0.0.0-20211019181941-9d821ace8654 // indirect
	golang.org/x/text v0.3.7 // indirect
)

replace et-protobuf3 => ../et-protobuf3

replace et-config => ../et-config
