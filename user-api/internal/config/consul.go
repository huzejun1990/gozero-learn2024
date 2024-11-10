// @Author huzejun 2024/11/11 5:23:00
package config

import (
	"bytes"
	"encoding/json"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/core/configcenter/subscriber"
	"sync"
)

type ConsulSubscriber struct {
	listeners []func()
	lock      sync.Mutex
	consulCli *consulApi.Client
	Path      string
}

func NewConsulSubscriber(address, path string) subscriber.Subscriber {
	config := &consulApi.Config{
		Address: address,
	}
	client, err := consulApi.NewClient(config)
	if err != nil {
		panic(err)
	}
	return &ConsulSubscriber{
		consulCli: client,
		Path:      path,
	}
}

func (m *ConsulSubscriber) AddListener(listener func()) error {
	m.lock.Lock()
	m.listeners = append(m.listeners, listener)
	m.lock.Unlock()
	return nil
}

func (m *ConsulSubscriber) Value() (string, error) {
	kvPair, _, err := m.consulCli.KV().Get(m.Path, nil)
	if err != nil {
		panic(err)
	}
	v := viper.New()
	v.SetConfigType("yaml")
	v.ReadConfig(bytes.NewBuffer(kvPair.Value))
	if err != nil {
		panic(err)
	}
	settings := v.AllSettings()
	marshal, _ := json.Marshal(settings)
	return string(marshal), nil
}
