package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hedon954/goapm"
	"github.com/hedon954/goapm/apm"

	"github.com/hedon954/goapm-example/protos"
	"github.com/hedon954/goapm-example/usrsvc/api"
	"github.com/hedon954/goapm-example/usrsvc/dao"
)

var Infra *goapm.Infra

func main() {
	Infra = goapm.NewInfra("usrsvc",
		goapm.WithMySQL("usrdb", "root:root@tcp(goapm-mysql:3306)/usrsvc?charset=utf8mb4&parseTime=True&loc=Local"),
		goapm.WithRedisV9("usrredis", "goapm-redis:6379", ""),
		goapm.WithAPM("goapm-otel-collector:4317"),
		goapm.WithMetrics(),
	)
	defer Infra.Stop()

	httpServer := apm.NewHTTPServer(":30012")
	grpcServer := apm.NewGrpcServer(":30002")
	protos.RegisterUserServiceServer(grpcServer.Server, &api.User{
		Dao: &dao.UserDao{DB: Infra.MySQL("usrdb"), RDB: Infra.RedisV9("usrredis")},
	})

	go httpServer.Start()
	go grpcServer.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-quit
}
