package servicePool

import (
	"loadBalancer/internal/balancingAlgorithm"
	"math"
	"math/rand"
	"sync/atomic"
)

func (obj *servicePool) DescribeNextServiceInstance() serviceInstance {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()

	if len(obj.pool) == 0 {
		return nil
	}
	switch obj.algorithm {
	case balancingAlgorithm.Random:
		return obj.pool[rand.Intn(len(obj.pool))]
	case balancingAlgorithm.LeastConnections:
		targetValue := int64(math.MaxInt64)
		targetValueIndex := 0
		targetCurrent := int64(0)

		for i, el := range obj.pool {
			targetCurrent = el.DescribeActiveConnections()
			if targetCurrent < targetValue {
				targetValue = targetCurrent
				targetValueIndex = i
			}
		}
		return obj.pool[targetValueIndex]
	case balancingAlgorithm.AverageResponseTime:
		targetValue := uint64(math.MaxUint64)
		targetValueIndex := 0
		targetCurrent := uint64(0)

		for i, el := range obj.pool {
			targetCurrent = el.DescribeAverageResponseTime()
			if targetCurrent < targetValue {
				targetValue = targetCurrent
				targetValueIndex = i
			}
		}
		return obj.pool[targetValueIndex]
	default:
		index := atomic.AddUint64(&obj.lastIndex, 1)
		return obj.pool[index%uint64(len(obj.pool))]
	}
}
