## Based on [exporttools](https://github.com/Zumata/exporttools)
Building blocks for quickly creating custom prometheus exporters.

Deprecated. see godoc

#### Getting Started

##### 1. Bootstrap your Exporter by creating a struct and embedding `*BaseExporter`. Your Exporter must satisfy the `Exporter` interface.
```
type Exporter interface {

	// to be implemented by custom exporter
	Setup() error
	Close() error

	Process()

	// satisfy via GenericCollect & GenericDescribe, or custom implementation
	prometheus.Collector
}
```

##### 2. Implement methods the `Setup()` and `Close()` methods required by the Exporter interface to create/destroy infrastructure connections.

##### 3. For each group of metrics to be collected by the Exporter, satisfy the `MetricCollector` interface then add to your Exporter via `AddCollector()` provided by the `BaseExporter`, within your Exporter's `Setup()` method.
```
type MetricCollector interface {
	Collect() ([]*Metric, error)
}
```

##### 4. Implement `prometheus.Collector` using the helpers:
```
func (e *customExporter) Describe(ch chan<- *prometheus.Desc) {
	exporttools.GenericDescribe(e.BaseExporter, ch)
}

func (e *customExporter) Collect(ch chan<- prometheus.Metric) {
	exporttools.GenericCollect(e.BaseExporter, ch)
}
```

##### 5. Calling `Register(exporter)` will enable collection for all metric groups.
```
func main()
  exporter := postgres.NewCustomExporter()
  err := exportertools.Register(exporter)
  if err != nil {
    log.Fatal(err)
  }
}
```
