package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/hedon954/goapm"
	"github.com/hedon954/goapm/apm"

	"github.com/hedon954/goapm-example/ordersvc/api"
	"github.com/hedon954/goapm-example/ordersvc/grpcclient"
	"github.com/hedon954/goapm-example/ordersvc/metric"
	"github.com/hedon954/goapm-example/protos"
)

var Infra *goapm.Infra

func main() {
	// init infra
	Infra = goapm.NewInfra("ordersvc",
		goapm.WithMySQL("orderdb", "root:root@tcp(goapm-mysql:3306)/ordersvc?charset=utf8mb4&parseTime=True&loc=Local"),
		goapm.WithAPM("goapm-otel-collector:4317"),
		goapm.WithMetrics(metric.All()...),
		goapm.WithAutoPProf(&apm.AutoPProfOpt{
			EnableCPU:       true,
			EnableMem:       true,
			EnableGoroutine: true,
		}),
	)

	// init grpc clients
	uc, err := apm.NewGrpcClient("goapm-ubuntu:30002", "usrsvc")
	if err != nil {
		panic(err)
	}
	grpcclient.UserClient = protos.NewUserServiceClient(uc)
	sc, err := apm.NewGrpcClient("goapm-ubuntu:30003", "skusvc")
	if err != nil {
		panic(err)
	}
	grpcclient.SkuClient = protos.NewSkuServiceClient(sc)

	// init http server
	hs := apm.NewHTTPServer(":30001")
	hs.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK"))
	})
	orderApi := &api.Order{
		Tracer: Infra.Tracer,
		DB:     Infra.MySQL("orderdb"),
	}
	hs.HandleFunc("/order/add", orderApi.Add)

	// start all servers
	go hs.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-quit
}
