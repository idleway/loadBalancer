package usecase

import (
	"context"
	"loadBalancer/internal/config"
)

func (obj *adminUsecase) DescribeConfig(ctx context.Context) config.ServicePool {
	return obj.servicePool.DescribeCurrentConfig(ctx)
}
func (obj *adminUsecase) StoreConfig(ctx context.Context, newConfig config.ServicePool) (err error) {
	return obj.servicePool.StoreNewConfig(ctx, newConfig)
}
