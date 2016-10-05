package exportertools

import(
    "github.com/prometheus/common/log"
    "github.com/prometheus/client_golang/prometheus"
    "time"
)

/*
BaseExporter provides convinient caching and metrics processing.
BaseExporter should be embedded into custom exporter struct
*/
type BaseExporter struct {
    Control chan bool
    Name    string
    Cache   *Cache
    MetricCollectors []MetricCollector
    Labels  map[string]string
}

func NewBaseExporter(name string, ttl int, labels map[string]string) *BaseExporter {
    e := BaseExporter{
        Control:  make(chan bool, 1),
        Name:     name,
        Cache:    NewCache(ttl),
        MetricCollectors:  []MetricCollector{},
        Labels:   labels,
    }
    return &e
}

func (e *BaseExporter) AddCollector(m MetricCollector) {
    e.MetricCollectors = append(e.MetricCollectors, m)
}

// Process runs metric collection in async mode
func (e *BaseExporter) Process() {
    defer log.Debugf("Instance terminated. Bye-bye from %v", e.Name, e)

    var channels []chan bool

    for _, _ = range e.MetricCollectors {
        channels = append(channels, make(chan bool))
    }

    for id, mc := range e.MetricCollectors {
        log.Debugf("Started Process() for %v", e.Name)
        ticker := time.NewTicker(e.Cache.TTL)
        go func(done chan bool, m MetricCollector) {
            defer close(done)
            for {
                select {
                case <-ticker.C:
                    go func() {
                        log.Debugf("Call Collect() for %v", e.Name)
                        metrics, err := m.Collect()
                        if err != nil {
                            log.Errorf("Error occured during Collect() for %v. error: %v", e.Name, err)
                        }
                        for idx := range metrics {
                            e.Cache.Set(metrics[idx])
                        }
                    }()
                case <-done:
                    log.Debugf("Exit from child of Process() %v", e.Name)
                    return
                }
            }
        }(channels[id], mc)
    }
    msg := <- e.Control
    go broadcastMessage(msg, channels)
}

func (e *BaseExporter) Close() error {
    defer close(e.Control)

    e.Control<- true
    log.Debugf("Stop processing metric for %v", e.Labels)
    err := Unregister(e)
    return err
}


func (e *BaseExporter) Setup() error {
    log.Debugf("Default Setup() func is used. It does nothing")
    return nil
}

// Satisfies Exporter interface and calls GenericDescribe which works with cache
func (e *BaseExporter) Describe(ch chan<- *prometheus.Desc) {
    GenericDescribe(e, ch)
}

// Satisfies Exporter interface and calls GenericCollect which works with cache
func (e *BaseExporter) Collect(ch chan<- prometheus.Metric) {
    GenericCollect(e, ch)
}

func broadcastMessage(msg bool, channels []chan bool) {
    for _, ch := range channels {
        ch <- msg
    }
}
