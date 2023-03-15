package serviceInstance

import (
	"net/url"
	"sync/atomic"
)

func (obj *serviceInstance) DescribeActiveConnections() int64 {
	return atomic.LoadInt64(&obj.stats.activeConnections)
}

func (obj *serviceInstance) DescribeAverageResponseTime() uint64 {
	return obj.stats.averageResponseTime.DescribeStatistics()
}

func (obj *serviceInstance) DescribeURL() *url.URL {
	return obj.url
}
