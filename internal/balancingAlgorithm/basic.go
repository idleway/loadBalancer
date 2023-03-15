package balancingAlgorithm

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type BalancingAlgorithm string

func (obj *BalancingAlgorithm) UnmarshalJSON(value []byte) (err error) {
	*obj, err = ValidateBalancingAlgorithm(string(value)[1 : len(string(value))-1])
	return
}
func (obj *BalancingAlgorithm) UnmarshalYAML(value *yaml.Node) (err error) {
	*obj, err = ValidateBalancingAlgorithm(value.Value)
	return
}

const (
	RoundRobin          BalancingAlgorithm = "roundRobin"
	Random              BalancingAlgorithm = "random"
	LeastConnections    BalancingAlgorithm = "leastConnections"
	AverageResponseTime BalancingAlgorithm = "averageResponseTime"
)

func ValidateBalancingAlgorithm(input string) (output BalancingAlgorithm, err error) {
	switch input {
	case string(RoundRobin):
		output = RoundRobin
	case string(Random):
		output = Random
	case string(LeastConnections):
		output = LeastConnections
	case string(AverageResponseTime):
		output = AverageResponseTime
	default:
		err = fmt.Errorf("unknown BalancingAlgorithm")
	}
	return
}
