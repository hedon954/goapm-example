package main

import (
	"github.com/hedon954/goapm"
	"github.com/hedon954/goapm-example/protos"
	"github.com/hedon954/goapm-example/skusvc/api"
)

func main() {
	dogapm.Infra.Init(
		dogapm.WithMySQL("root:root@tcp(apm-mysql:3306)/skusvc?charset=utf8mb4&parseTime=True&loc=Local"),
		dogapm.WithEnableAPM("apm-otel-collector:4317", "/logs", 10),
		dogapm.WithMetric(),
	)

	dogapm.NewHTTPServer(":30013")
	gs := dogapm.NewGrpcServer(":30003")
	protos.RegisterSkuServiceServer(gs, &api.SkuService{})

	dogapm.EndPoint.Start()
}
