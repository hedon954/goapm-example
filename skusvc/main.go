package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hedon954/goapm"
	"github.com/hedon954/goapm/apm"

	"github.com/hedon954/goapm-example/protos"
	"github.com/hedon954/goapm-example/skusvc/api"
	"github.com/hedon954/goapm-example/skusvc/dao"
)

var Infra *goapm.Infra

func main() {
	Infra = goapm.NewInfra("skusvc",
		goapm.WithMySQL("skudb", "root:root@tcp(goapm-mysql:3306)/skusvc?charset=utf8mb4&parseTime=True&loc=Local"),
		goapm.WithAPM("goapm-otel-collector:4317"),
		goapm.WithMetrics(),
	)

	httpServer := apm.NewHTTPServer(":30013")
	grpcServer := apm.NewGrpcServer(":30003")
	protos.RegisterSkuServiceServer(grpcServer, &api.SkuService{
		Dao: &dao.SkuDao{DB: Infra.MySQL("skudb"), RDB: Infra.RedisV9("skuredis")},
	})

	go httpServer.Start()
	go grpcServer.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-quit
}
