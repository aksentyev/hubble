package exportertools_mock

import "github.com/aksentyev/hubble/exportertools"

// Mock of BackendAdapter interface
type Collector struct {
}

func NewMockCollector() *Collector {
    mock := Collector{}
    return &mock
}

func (c *Collector) Collect() ([]*exportertools.Metric, error) {
    var mock []*exportertools.Metric

    metric := exportertools.Metric {
        "test_collector_metric",
        "test descr",
        exportertools.StringToType("Gauge"),
        23,
        map[string]string{},
    }

    mock = append(mock, &metric)

    return mock, nil
}
