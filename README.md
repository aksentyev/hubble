# Hubble service discovery

Deprecated. see godoc

Hubble is the alternative service discovery for Prometheus. Consul is the only one supported backend.

**Usage:**

```go
import (
	"github.com/aksentyev/hubble/hubble"
	"github.com/aksentyev/hubble/consul"
)

params := &consul.Config{
  Address: "consul.service.consul:8500",
  Datacenter: "staging",
}

exporterName := "postgres"

client, _ := consul.New(params)
kv := consul.NewKV(client)
h := hubble.New(client, kv, exporterName)

svcs := func() (list []*hubble.ServiceAtomic) {
	for _, svc := range h.Services(){
		for _, el := range svc.MakeAtomic(nil) {
			list = append(list, el)
		}
	}
	return list
}

c = h.NewCache(60 * time.Second, svcs)
```

`h.NewCache(60 * time.Second, svcs)`
The first parameter is TTL for the cache, the second is callback function that should return []*ServiceAtomic list.

ServiceAtomic struct was designed as a prepared unit to be used as a parameter for an exporter:

***Example:***

```go
s = svc.MakeAtomic(additional_labels_map) //map[string]string
```

Consul KV keeps json object contains paramethers could be used to automatically enable collecting metric for the service.

Path mask: `monitoring/$service/$exporter_name`

E.g it can be useful for Alertmanager rules

***Example:***

```json
{
    "notifiable": true,
	"extra_labels": {"env": "staging"},
    "exporter_options": {"password": "123456"}
}
```
