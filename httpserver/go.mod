module httpserver

go 1.16

require (
	et-protobuf3 v0.0.0
	github.com/gin-gonic/gin v1.8.1
	github.com/golang-jwt/jwt v3.2.2+incompatible
	google.golang.org/grpc v1.48.0
)

replace et-protobuf3 => ../et-protobuf3
