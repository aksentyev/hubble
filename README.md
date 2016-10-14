# Hubble service discovery

[![GoDoc](https://godoc.org/github.com/aksentyev/hubble?status.svg)](https://godoc.org/github.com/aksentyev/hubble) [![Build Status](https://travis-ci.org/aksentyev/hubble.svg?branch=master)](https://travis-ci.org/aksentyev/hubble) [![codecov](https://codecov.io/gh/aksentyev/hubble/branch/master/graph/badge.svg)](https://codecov.io/gh/aksentyev/hubble)

Hubble is the alternative service discovery for Prometheus. Consul is the only one supported backend.

Example usage:

```go
config := consul.DefaultConfig()
config.Address = *consulURL
config.Datacenter = *consulDC

client, err := consul.New(config)
if err != nil {
    panic(err)
}

kv := consul.NewKV(client)
h := hubble.New(client, kv, *consulTag)

// Filter services by any criteria, e.g. by consul tag.
filterCB := func(list []*hubble.Service) []*hubble.Service {
    var servicesForMonitoring []*hubble.Service
    for _, svc := range list {
        if util.IncludesStr(svc.Tags, *consulTag) {
            servicesForMonitoring = append(servicesForMonitoring, svc)
        }
    }
    return servicesForMonitoring
}

list := h.Services(filterCB)
```

Consul KV keeps json object contains parameters could be used to automatically enable collecting metric for the service.

Path mask: `monitoring/$service/$exporter_name`

E.g it can be useful for Alertmanager rules or any other parameters.

***Example:***

```json
{
    "notifiable": true,
    "extra_labels": {"env": "staging"},
    "exporter_options": {"password": "123456"}
}
```

## Exportertools

Quickly creating custom thread-safe prometheus exporters with SD and cache.
