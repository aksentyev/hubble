package consul

import (
    "github.com/hashicorp/consul/api"
)

type Consul struct {
    *api.Client
}

// Consul Config representation
type Config api.Config

type ConsulKV struct {
    *api.KV
}
