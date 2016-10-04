package exportertools

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/common/log"

    "strings"
    "math"
    "strconv"
)

type Metric struct {
    Name        string
    Description string
    Type        MetricType
    Value       interface{}
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

func (m *Metric) PromValue() float64 {
    switch v := m.Value.(type) {
    case int64:
        return float64(v)
    case float64:
        return v
    case time.Time:
        return float64(v.Unix())
    case []byte:
        // Try and convert to string and then parse to a float64
        strV := string(v)
        result, err := strconv.ParseFloat(strV, 64)
        if err != nil {
            return math.NaN()
        }
        return result
    case string:
        result, err := strconv.ParseFloat(v, 64)
        if err != nil {
            return math.NaN()
        }
        return result
    case nil:
        return math.NaN()
    default:
        return math.NaN()
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
