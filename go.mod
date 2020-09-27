module github.com/linger1216/go-utils

go 1.14

replace (
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)

require (
	github.com/ClickHouse/clickhouse-go v1.4.3
	github.com/araddon/dateparse v0.0.0-20200409225146-d820a6159ab1
	github.com/bwmarrin/snowflake v0.3.0
	github.com/coreos/etcd v3.3.22+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dgraph-io/ristretto v0.0.3
	github.com/gogo/protobuf v1.3.1
	github.com/golang/geo v0.0.0-20200730024412-e86565bf3f35
	github.com/google/uuid v1.1.2 // indirect
	github.com/gookit/config/v2 v2.0.17
	github.com/jmoiron/sqlx v1.2.0
	github.com/json-iterator/go v1.1.10
	github.com/lib/pq v1.0.0
	go.uber.org/zap v1.16.0
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/genproto v0.0.0-20200925023002-c2d885f95484 // indirect
)
