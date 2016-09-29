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
func (c *Consul) GetAll() (services []*hubble.Service, err error) {
    var wg sync.WaitGroup
    threads := make(chan bool, 10)
    servicesCh := make(chan *hubble.Service)
    errCh := make(chan error)

    list, err := c.getList()
    if err != nil {
        log.Errorf("Consul. Get services list failed: %+v", err)
        return services, err
    }

    wg.Add(len(list))

    go func() {
        defer close(threads)
        defer close(servicesCh)
        for _, name := range list {
            threads<- true
            go c.getService(servicesCh, threads, errCh, &wg, name)
        }
        wg.Wait()
    }()

    select {
    case err = <-errCh:
        close(errCh)
        return services, err
    default:
        for svc := range servicesCh {
            services = append(services, svc)
        }
    }
    log.Debugf("All services were fetched")

    return services, err
}

func (c *Consul) getService(servicesCh chan *hubble.Service, threads chan bool, errCh chan error, wg *sync.WaitGroup, name string) {
    defer func() { <-threads }()
    defer wg.Done()

    svc, err := c.fetch(name)
    if err != nil {
        errCh <- err
        log.Errorf("Consul: Get svc %v failed", svc.Name)
    }
    servicesCh<- svc
}


func (c *Consul) fetch(name string) (svc *hubble.Service, err error) {
    svc = hubble.DefaultService()
    svc.Name = name

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
    return svc, err
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
