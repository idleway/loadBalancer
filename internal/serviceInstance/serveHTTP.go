package serviceInstance

import (
	"net/http"
	"sync/atomic"
	"time"
)

func (obj *serviceInstance) ServeHTTP(rw http.ResponseWriter, req *http.Request, cacheResponse bool) {
	atomic.AddInt64(&obj.stats.activeConnections, 1)
	defer atomic.AddInt64(&obj.stats.activeConnections, -1)
	start := time.Now()
	if cacheResponse {
		obj.reverseProxyWithCache.ServeHTTP(rw, req)
	} else {
		obj.reverseProxyNoCache.ServeHTTP(rw, req)
	}

	obj.stats.averageResponseTime.UpdateStatistics(uint64(time.Since(start).Milliseconds()))
}
