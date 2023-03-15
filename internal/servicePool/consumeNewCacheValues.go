package servicePool

import (
	"context"
	"loadBalancer/internal/bytesPool"
	cacheModel "loadBalancer/internal/cache"
)

func (obj *servicePool) consumeNewCacheValues(ctx context.Context) {
	var newCacheValue cacheModel.StoreInCache
	for {
		select {
		case <-ctx.Done():
			return
		case newCacheValue = <-obj.newCacheValues:
			obj.cache.Store(newCacheValue.Key, newCacheValue.RespBody)

			newCacheValue.RespBody.Reset()
			bytesPool.Pool.Put(newCacheValue.RespBody)
		}
	}
}
