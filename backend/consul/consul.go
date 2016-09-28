package consul

import (
	"github.com/aksentyev/consul/api"
	"github.com/aksentyev/hubble/hubble"
	"sync"
	"strconv"
	"encoding/json"
    "github.com/prometheus/common/log"
)
// Convert api.DefaultConfig to Config
func DefaultConfig() Config {
    return Config(*api.DefaultConfig())
}

// Creates new consul api instance according to the configuration was received
func New(src Config) (consul *Consul, err error) {
    converted := api.Config(src)
	client, err := api.NewClient(&converted)
	consul = &Consul{client}
	return consul, err
}

// Returns new instance of the KV
func NewKV(consul *Consul) *ConsulKV {
	kv := &ConsulKV{consul.KV()}
	return kv
}

// Return Services list
func (c *Consul) GetAll() []*hubble.Service {
	var (
		wg sync.WaitGroup
		mutex sync.Mutex
	)
	threads := make(chan bool, 10)

	list, err := c.getList()
	if err != nil {
		log.Errorf("Consul. Get services list failed: %+v", err)
	}
	services := []*hubble.Service{}
	defer close(threads)

	wg.Add(len(list))

	go func() {
		for _, name := range list {
            threads<- true
			go func(name string) {
                defer func() { <-threads }()
				defer wg.Done()
				defer mutex.Unlock()

                svc := hubble.DefaultService()
                svc.Name = name

				err := c.get(svc)
				if err != nil {
					log.Errorf("Consul: Get svc %v failed", svc.Name)
				}
				mutex.Lock()
				services = append(services, svc)
			}(name)
		}
	}()

	wg.Wait()
	return services
}


func (c *Consul) get(svc *hubble.Service) (err error) {
	res, _, err := c.Catalog().Service(svc.Name, "", nil)

	for _, s := range res {
		if s.ServiceAddress == "" {
			svc.Addresses[s.Node] = s.Address
		} else {
			svc.Addresses[s.Node] = s.ServiceAddress
		}
		svc.Port = strconv.Itoa(s.ServicePort)
		svc.Tags = s.ServiceTags
        svc.ModifyIndex = s.ModifyIndex
	}
	return err
}

func (c *Consul) getList() (list []string, err error) {
	res, _, err := c.Catalog().Services(nil)
	for k, _ := range res {
		list = append(list, k)
	}
	return list, err
}

// Returns values from the KV
func (c *ConsulKV) GetParams(key string) (p *hubble.ServiceParams, err error) {
	var res *api.KVPair
	p = &hubble.ServiceParams{}

	if res, _, err = c.Get(key, nil); res != nil {
		err = json.Unmarshal(res.Value, p)
        p.ModifyIndex = res.ModifyIndex
	}
	return p, err
}
