[![GoDoc](https://godoc.org/github.com/aksentyev/hubble/exportertools?status.svg)](https://godoc.org/github.com/aksentyev/hubble/exportertools)

#### Getting Started

##### 1. Bootstrap your Exporter by creating a struct and embedding `*BaseExporter`. Your Exporter must satisfy the `Exporter` interface.
```go
type Exporter interface {

	// to be implemented by custom exporter
	Setup() error
	Close() error

	Process()

	// satisfy via GenericCollect & GenericDescribe, or custom implementation
	prometheus.Collector
}
```

Note that `(\*BaseExporter).Close()` should be called from custom `Close()` for gracefully shutdown exporter's `Process()` goroutines.

##### 2. Implement methods the `Setup()` and `Close()` methods required by the Exporter interface to create/destroy infrastructure connections.

##### 3. For each group of metrics to be collected by the Exporter, satisfy the `MetricCollector` interface then add to your Exporter via `AddCollector()` provided by the `BaseExporter`, within your Exporter's `Setup()` method.

```go
type MetricCollector interface {
	Collect() ([]*Metric, error)
}
```

##### 4. Implement `prometheus.Collector` using the helpers:

```go
func (e *customExporter) Describe(ch chan<- *prometheus.Desc) {
	exporttools.GenericDescribe(e.BaseExporter, ch)
}

func (e *customExporter) Collect(ch chan<- prometheus.Metric) {
	exporttools.GenericCollect(e.BaseExporter, ch)
}
```

##### 5. Calling `Register(exporter)` will enable collection for all metric groups.
```go
func main()
  exporter := postgres.NewCustomExporter()
  err := exportertools.Register(exporter)
  if err != nil {
    log.Fatal(err)
  }
}
```

##### 6. *Dispatcher* component gets services, sends new services to be Registered to the channel and services to be Unregistered to another one.

- Send callback function to the Dispatcher instance. Function should return actual services list.

```go
cb := func() (list []*hubble.ServiceAtomic, err error) {
    services, err := h.Services(filterCB)
    if err != nil {
        return list, err
    }
    for _, svc := range services {
        for _, el := range svc.MakeAtomic(nil) {
            list = append(list, el)
        }
    }
    return list, err
}
```

- Set interval between *cb* executions, then run Dispatcher in background

```go
d = hubble.NewDispatcher(*updateInterval)
go d.Run(cb)
```

- Listen for `ToRegister` and `ToUnregister` channels for new objects. Register and Unregister theirs exporters.

```go
func listenAndRegister() {
    pgMetricsParsed := exporter.AddFromFile(*queriesPath)

    for svc := range d.ToRegister {
        if len(svc.ExporterOptions) > 1 {
            config := exporter.Config{
                DSN:             util.PgConnURL(svc),
                Labels:          svc.ExtraLabels,
                ExporterOptions: svc.ExporterOptions,
                CacheTTL:        *scrapeInterval,
                PgMetrics:       pgMetricsParsed,
            }
            exp, err := exporter.CreateAndRegister(&config)
            if err == nil {
                d.Register(svc, exp)
                log.Infof("Registered %v %v", svc.Name, svc.Address)
            } else {
                log.Warnf("Register was failed for service %v %v %v", svc.Name, svc.Address, err)
                exp.Close()
            }
        }
    }
}

func listenAndUnregister() {
    for m := range d.ToUnregister {
        for h, svc := range m {
            exporter := d.Exporters[h].(*exporter.PostgresExporter)
            err := exporter.Close()
            if err != nil {
                log.Warnf("Unregister() for %v %v returned %v:", svc.Name, svc.Address, err)
            } else {
                log.Infof("Unregister service %v %v", svc.Name, svc.Address)
            }
            d.UnregisterWithHash(h)
        }
    }
}

```

##### Inspired by [exporttools](https://github.com/Zumata/exporttools)
