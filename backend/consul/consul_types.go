package consul

import (
    "github.com/aksentyev/consul/api"
)

type Consul struct {
    *api.Client
}

// Consul Config representation
type Config api.Config

type ConsulKV struct {
    *api.KV
}
