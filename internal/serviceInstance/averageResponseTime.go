package serviceInstance

import (
	"context"
	"sync/atomic"
	"time"
)

const updateInterval = time.Second
const slidingWindowInterval = time.Minute * 5

type quantum struct {
	sumMs uint64
	count uint64
}
type averageResponseTime struct {
	data    []quantum
	pointer uint

	total quantum
	avgMs uint64

	collectorChan chan uint64
}

func NewAverageResponseTime(ctx context.Context) *averageResponseTime {
	obj := averageResponseTime{
		data:          make([]quantum, int(slidingWindowInterval.Seconds())),
		pointer:       0,
		collectorChan: make(chan uint64), //TODO нужен ли буферизованный?
	}
	go obj.consume(ctx, updateInterval)
	return &obj
}

func (obj *averageResponseTime) consume(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var newValue uint64
	for {
		select {
		case <-ticker.C:
			obj.movePointer()
		default:

		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			obj.movePointer()
		case newValue = <-obj.collectorChan:
			obj.data[obj.pointer].count++
			obj.data[obj.pointer].sumMs += newValue
		}
	}
}
func (obj *averageResponseTime) movePointer() {
	obj.total.sumMs += obj.data[obj.pointer].sumMs
	obj.total.count += obj.data[obj.pointer].count

	obj.pointer = (obj.pointer + 1) % uint(len(obj.data))

	obj.total.sumMs -= obj.data[obj.pointer].sumMs
	obj.total.count -= obj.data[obj.pointer].count

	obj.data[obj.pointer].sumMs = 0
	obj.data[obj.pointer].count = 0

	if obj.total.count == 0 {
		atomic.StoreUint64(&obj.avgMs, 0)
	} else {
		atomic.StoreUint64(&obj.avgMs, obj.total.sumMs/obj.total.count)
	}
}

func (obj *averageResponseTime) DescribeStatistics() (responseTimeMs uint64) {
	return atomic.LoadUint64(&obj.avgMs)
}
func (obj *averageResponseTime) UpdateStatistics(responseTimeMs uint64) {
	obj.collectorChan <- responseTimeMs
}
