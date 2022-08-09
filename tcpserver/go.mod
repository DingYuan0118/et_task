module tcpserver

go 1.16

require (
	et-config v0.0.0
	et-protobuf3 v0.0.0
	github.com/go-micro/plugins/v4/registry/etcd v1.1.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gomodule/redigo v1.8.9
	go-micro.dev/v4 v4.8.0
	go.etcd.io/etcd/client/v3 v3.5.4 // indirect
	go.uber.org/zap v1.17.0
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e
	xorm.io/xorm v1.3.1
)

replace et-protobuf3 => ../et-protobuf3

replace et-config => ../et-config
