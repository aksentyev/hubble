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

##### 3. For each group of metrics to be collected by the Exporter, satisfy the `MetricGroup` interface then add to your Exporter via `AddMetric()` provided by the `BaseExporter`, within your Exporter's `Setup()` method.
```
type Metric interface {
	Collect() ([]*Metric, error)
}
```

##### 4. `Metrics` handle both counters and gauges
```
type Metric struct {
	Name        string
	Description string
	Type        metricType
	Value       int64
	LabelKeys   []string
	LabelVals   []string
}

type metricType int

const (
	Counter   metricType = iota
    Gauge     metricType = iota
    Histogram metricType = iota
    Summary   metricType = iota
)
```

##### 5. Implement `prometheus.Collector` using the helpers:
```
func (e *customExporter) Describe(ch chan<- *prometheus.Desc) {
	exporttools.GenericDescribe(e.BaseExporter, ch)
}

func (e *customExporter) Collect(ch chan<- prometheus.Metric) {
	exporttools.GenericCollect(e.BaseExporter, ch)
}
```

##### 6. Calling `Export(exporter)` will enable collection for all metric groups.
```
func main()
  exporter := postgres.NewCustomExporter()
  err := exportertools.Export(exporter)
  if err != nil {
    log.Fatal(err)
  }
}
```
