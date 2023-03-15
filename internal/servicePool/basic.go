package servicePool

import (
	"bytes"
	"context"
	"loadBalancer/internal/balancingAlgorithm"
	cacheModel "loadBalancer/internal/cache"
	"loadBalancer/internal/config"
	serviceInstanceModel "loadBalancer/internal/serviceInstance"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

type serviceInstance interface {
	ServeHTTP(rw http.ResponseWriter, req *http.Request, cacheResponse bool)

	DescribeAverageResponseTime() uint64
	DescribeActiveConnections() int64
	DescribeURL() *url.URL
}
type cache interface {
	Describe(key uint64) (value []byte, exist bool)
	Store(key uint64, value *bytes.Buffer)
}
type launchedServices struct {
	serviceInstance
	ctx    context.Context
	cancel context.CancelFunc
}
type servicePool struct {
	ctx context.Context

	pool  []launchedServices
	mutex sync.RWMutex

	cache          cache
	newCacheValues chan cacheModel.StoreInCache
	cacheEnabled   atomic.Bool

	algorithm balancingAlgorithm.BalancingAlgorithm
	lastIndex uint64

	currentConfig      config.ServicePool
	currentConfigMutex sync.RWMutex
}

func NewServicePool(ctx context.Context, config config.ServicePool) *servicePool {
	obj := servicePool{
		ctx: ctx,

		pool:      make([]launchedServices, len(config.ServicesPool)),
		algorithm: config.BalancingAlgorithm,

		cache:          cacheModel.NewCache(),
		newCacheValues: make(chan cacheModel.StoreInCache),

		currentConfig: config,
	}

	obj.cacheEnabled.Store(config.CacheEnabled)
	for i, el := range config.ServicesPool {
		obj.pool[i].ctx, obj.pool[i].cancel = context.WithCancel(ctx)
		obj.pool[i].serviceInstance = serviceInstanceModel.NewServiceInstance(obj.pool[i].ctx, el.URL, obj.newCacheValues)
	}

	rand.Seed(time.Now().Unix())
	go obj.consumeNewCacheValues(ctx)
	return &obj
}
