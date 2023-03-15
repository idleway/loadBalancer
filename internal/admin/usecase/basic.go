package usecase

import (
	"context"
	"loadBalancer/internal/config"
)

type servicePool interface {
	DescribeCurrentConfig(ctx context.Context) config.ServicePool
	StoreNewConfig(ctx context.Context, newConfig config.ServicePool) (err error)
}

type adminUsecase struct {
	servicePool servicePool
}

func NewAdminUsecase(servicePool servicePool) *adminUsecase {
	return &adminUsecase{servicePool: servicePool}
}
