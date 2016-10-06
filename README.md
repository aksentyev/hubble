# Hubble service discovery

[![GoDoc](https://godoc.org/github.com/aksentyev/hubble?status.svg)](https://godoc.org/github.com/aksentyev/hubble)

Hubble is the alternative service discovery for Prometheus. Consul is the only one supported backend.

Example usage see in [postgres_exporter](https://github.com/aksentyev/postgres_exporter) code.

Consul KV keeps json object contains parameters could be used to automatically enable collecting metric for the service.

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
