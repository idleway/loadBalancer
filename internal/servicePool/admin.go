package servicePool

import (
	"context"
	"loadBalancer/internal/config"
	serviceInstanceModel "loadBalancer/internal/serviceInstance"
)

func (obj *servicePool) DescribeCurrentConfig(ctx context.Context) config.ServicePool {
	obj.currentConfigMutex.RLock()
	defer obj.currentConfigMutex.RUnlock()
	return obj.currentConfig
}
func (obj *servicePool) StoreNewConfig(ctx context.Context, newConfig config.ServicePool) (err error) {
	obj.currentConfigMutex.Lock()
	defer obj.currentConfigMutex.Unlock()

	if obj.currentConfig.Equal(newConfig) {
		return
	}

	obj.mutex.Lock()
	defer obj.mutex.Unlock()

	obj.cacheEnabled.Store(newConfig.CacheEnabled)
	obj.algorithm = newConfig.BalancingAlgorithm

	newConfigServicesMap := make(map[string]int)
	for i := range newConfig.ServicesPool {
		newConfigServicesMap[newConfig.ServicesPool[i].String()] = i
	}

	for i := 0; ; {
		_, exist := newConfigServicesMap[obj.pool[i].serviceInstance.DescribeURL().String()]
		if exist {
			delete(newConfigServicesMap, obj.pool[i].serviceInstance.DescribeURL().String())
			i++
		} else {
			obj.pool[i] = obj.pool[len(obj.pool)-1]
			obj.pool = obj.pool[:len(obj.pool)-1]
		}
		if i >= len(obj.pool) {
			break
		}
	}
	for _, value := range newConfigServicesMap {
		newService := launchedServices{serviceInstance: nil}
		newService.ctx, newService.cancel = context.WithCancel(obj.ctx)
		newService.serviceInstance = serviceInstanceModel.NewServiceInstance(newService.ctx, newConfig.ServicesPool[value].URL, obj.newCacheValues)
		obj.pool = append(obj.pool, newService)
	}
	obj.lastIndex = 0
	obj.currentConfig = newConfig

	return nil
}
