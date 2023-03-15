package config

import (
	"bytes"
	"encoding/json"
	"gopkg.in/yaml.v3"
	"loadBalancer/internal/balancingAlgorithm"
	"net/url"
)

var cfg *Config

func GetConfigInstance() Config {
	if cfg != nil {
		return *cfg
	}
	return Config{}
}
func ReadConfigYML(configBytes []byte) error {
	decoder := yaml.NewDecoder(bytes.NewReader(configBytes))
	if err := decoder.Decode(&cfg); err != nil {
		return err
	}
	return nil
}

type ServiceInstance struct {
	*url.URL
}

func (obj *ServiceInstance) MarshalJSON() ([]byte, error) {
	return json.Marshal(obj.String())
}
func (obj *ServiceInstance) UnmarshalJSON(value []byte) (err error) {
	obj.URL, err = url.Parse(string(value)[1 : len(string(value))-1])
	return
}
func (obj *ServiceInstance) UnmarshalYAML(value *yaml.Node) (err error) {
	obj.URL, err = url.Parse(value.Value)
	return
}

type ServicePool struct {
	CacheEnabled       bool                                  `yaml:"cacheEnabled" json:"cacheEnabled"`
	BalancingAlgorithm balancingAlgorithm.BalancingAlgorithm `yaml:"balancingAlgorithm" json:"balancingAlgorithm"`
	ServicesPool       []ServiceInstance                     `yaml:"servicesPool" json:"servicesPool"`
}

func (obj *ServicePool) Equal(reference ServicePool) bool {
	if len(obj.ServicesPool) != len(reference.ServicesPool) ||
		obj.CacheEnabled != reference.CacheEnabled ||
		obj.BalancingAlgorithm != reference.BalancingAlgorithm {
		return false
	}
	temp := make(map[string]any)
	for i := range obj.ServicesPool {
		temp[obj.ServicesPool[i].String()] = nil
	}
	for i := range reference.ServicesPool {
		_, exist := temp[reference.ServicesPool[i].String()]
		if !exist {
			return false
		}
	}
	return true
}

type Project struct {
	Name      string `yaml:"name"`
	PprofPort string `yaml:"pprofPort"`
}

type Config struct {
	Project     Project     `yaml:"project"`
	ServicePool ServicePool `yaml:"servicePool"`
}
