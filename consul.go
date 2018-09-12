package consul

import (
	"log"
	consulapi "github.com/hashicorp/consul/api"
)

var serverAddress string
var consulClient *consulapi.Client

type ServiceRegistration struct {
	ID      string
	Name    string
	Tags    []string
	Port    int
	Address string
	Check   ServiceCheck
}
type ServiceCheck struct {
	Interval                       string
	Timeout                        string
	HTTP                           string
	DeregisterCriticalServiceAfter string
}

func CheckErr(err error) {
	if err != nil {
		log.Printf("[E]: %v", err)
	}
}

func Init(serverAddress string) {
	serverAddress = serverAddress

	config := consulapi.DefaultConfig()
	config.Address = serverAddress
	var err error
	consulClient, err = consulapi.NewClient(config)
	if err != nil {
		log.Fatal("consul client error : ", err)
	}
}
func RegisterServer(info ServiceRegistration) {

	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = info.ID
	registration.Name = info.Name
	registration.Port = info.Port
	registration.Tags = info.Tags
	registration.Address = info.Address
	registration.Check = &consulapi.AgentServiceCheck{
		HTTP:                           info.Check.HTTP,
		Timeout:                        info.Check.Timeout,
		Interval:                       info.Check.Interval,
		DeregisterCriticalServiceAfter: info.Check.DeregisterCriticalServiceAfter, //check失败后30秒删除本服务
	}
	var err error
	err = consulClient.Agent().ServiceRegister(registration)
	CheckErr(err)
}

func GetKV(key string) string {

	kv, _, err := consulClient.KV().Get(key, nil)
	CheckErr(err)
	if kv == nil {
		return ""
	} else {
		return string(kv.Value)
	}
}
