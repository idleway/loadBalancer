package serviceInstance

import (
	"context"
	"loadBalancer/internal/cache"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type averageResponseTimeI interface {
	DescribeStatistics() (responseTimeMs uint64)
	UpdateStatistics(responseTimeMs uint64)
}
type stats struct {
	activeConnections   int64
	averageResponseTime averageResponseTimeI
}

type serviceInstance struct {
	url   *url.URL
	stats stats

	reverseProxyNoCache   *httputil.ReverseProxy
	reverseProxyWithCache *httputil.ReverseProxy
}

func NewServiceInstance(ctx context.Context, url *url.URL, cacheChan chan cache.StoreInCache) *serviceInstance {
	obj := serviceInstance{
		url: url,
		stats: stats{
			activeConnections:   0,
			averageResponseTime: NewAverageResponseTime(ctx),
		},

		reverseProxyNoCache:   httputil.NewSingleHostReverseProxy(url),
		reverseProxyWithCache: httputil.NewSingleHostReverseProxy(url),
	}
	obj.reverseProxyWithCache.Transport = &transportWithCache{http.DefaultTransport, cacheChan}
	return &obj
}
