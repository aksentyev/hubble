package exportertools

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/common/log"

    "strings"
)

type Metric struct {
    Name        string
    Description string
    Type        MetricType
    Value       float64
    Labels      map[string]string
}

type MetricType int

type MetricCollector interface {
    Collect() ([]*Metric, error)
}

const (
    Counter MetricType = iota
    Gauge   MetricType = iota
    Untyped MetricType = iota
)

func (m *Metric) PromDescription(exporterName string) *prometheus.Desc {
    return prometheus.NewDesc(
        prometheus.BuildFQName("", exporterName, m.Name),
        m.Description,
        nil, prometheus.Labels(m.Labels),
    )
}

func (m *Metric) PromType() prometheus.ValueType {
    switch m.Type {
    case Counter:
        return prometheus.CounterValue
    case Gauge:
        return prometheus.GaugeValue
    default:
        return prometheus.UntypedValue
    }
}

func StringToType(s string) MetricType {
    switch strings.ToLower(s) {
    case "gauge":
        return Gauge
    case "counter":
        return Counter
    case "untyped":
        return Untyped
    default:
        log.Errorf("Undefined metric type: %v", s)
        return Untyped
    }
}
