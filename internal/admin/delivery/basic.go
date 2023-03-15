package delivery

import (
	"context"
	"loadBalancer/internal/config"
	"net/http"
)

type adminUsecase interface {
	DescribeConfig(ctx context.Context) config.ServicePool
	StoreConfig(ctx context.Context, newConfig config.ServicePool) (err error)
}

type adminHandler struct {
	adminUsecase adminUsecase
}

func NewAdminHandler(r *http.ServeMux, adminUsecase adminUsecase) {
	handler := &adminHandler{adminUsecase: adminUsecase}

	r.HandleFunc("/admin/describe_config", handler.describeConfig)
	r.HandleFunc("/admin/store_config", handler.storeConfig)
}
