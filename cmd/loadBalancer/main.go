package main

import (
	"context"
	"fmt"
	adminDelivery "loadBalancer/internal/admin/delivery"
	adminUsecase "loadBalancer/internal/admin/usecase"
	"loadBalancer/internal/config"
	"loadBalancer/internal/servicePool"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	// hey -z 15s -c 8 -q 50 -m POST http://localhost:9000/easy
	// hey -n 30 -c 8 -m POST http://localhost:9000/easy
	// hey -z 15s -c 4 -q 8 -m POST http://localhost:9000/easy
	ctx := context.Background()

	configBytes, err := os.ReadFile("config/config.yml")
	if err != nil {
		log.Fatalln("Failed read configuration file", err)
	}
	if err := config.ReadConfigYML(configBytes); err != nil {
		log.Fatalln("Failed init local configuration", err)
	}
	adminMux := http.ServeMux{}

	servicePoolInstance := servicePool.NewServicePool(ctx, config.GetConfigInstance().ServicePool)
	adminUsecaseInstance := adminUsecase.NewAdminUsecase(servicePoolInstance)
	adminDelivery.NewAdminHandler(&adminMux, adminUsecaseInstance)

	lbSrv := http.Server{
		Addr:    fmt.Sprintf(":%d", 9000),
		Handler: http.HandlerFunc(servicePoolInstance.LBFunc),
	}
	adminSrv := &http.Server{
		Addr:    fmt.Sprintf(":%d", 9100),
		Handler: &adminMux,
	}

	go lbSrv.ListenAndServe()
	go adminSrv.ListenAndServe()

	http.ListenAndServe(fmt.Sprintf("0.0.0.0:9200"), nil)
}
