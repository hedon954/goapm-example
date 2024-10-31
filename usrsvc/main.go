package main

import (
	"github.com/hedon954/goapm"
	"github.com/hedon954/goapm-example/protos"
	"github.com/hedon954/goapm-example/usrsvc/api"
)

func main() {
	dogapm.Infra.Init(
		dogapm.WithMySQL("root:root@tcp(apm-mysql:3306)/usrsvc?charset=utf8mb4&parseTime=True&loc=Local"),
		dogapm.WithRedis("apm-redis:6379"),
		dogapm.WithEnableAPM("apm-otel-collector:4317", "/logs", 10),
		dogapm.WithMetric(),
	)

	dogapm.NewHTTPServer(":30012")
	hs := dogapm.NewGrpcServer(":30002")
	protos.RegisterUserServiceServer(hs.Server, &api.User{})

	dogapm.EndPoint.Start()
}
