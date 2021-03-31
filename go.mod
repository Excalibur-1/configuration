module github.com/Excalibur-1/configuration

go 1.13

replace (
	github.com/coreos/bbolt v1.3.4 => go.etcd.io/bbolt v1.3.4
	google.golang.org/grpc v1.29.1 => google.golang.org/grpc v1.26.0
)

require (
	github.com/Excalibur-1/gutil v0.0.0-20210330073727-af5cc5f0f34a
	github.com/coreos/bbolt v1.3.4 // indirect
	github.com/coreos/etcd v3.3.25+incompatible // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/prometheus/client_golang v1.10.0 // indirect
	github.com/samuel/go-zookeeper v0.0.0-20201211165307-7117e9ea2414
	go.etcd.io/etcd v3.3.25+incompatible
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
	google.golang.org/grpc v1.29.1 // indirect
)
