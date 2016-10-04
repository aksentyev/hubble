package exportertools

import (
    "github.com/prometheus/common/log"
    "github.com/prometheus/client_golang/prometheus"
)

/*
GenericDescribe generates prom descriptor from cached metric names
*/
func GenericDescribe(be *BaseExporter, ch chan<- *prometheus.Desc) {
    fqName := prometheus.BuildFQName("", be.Name, "scrapes")
    labels := prometheus.Labels(be.Labels)
    log.Debugf("GenericDescribe: %+v %v", fqName, labels)
    ch <- prometheus.NewDesc(
        fqName,
        "times exporter has been scraped",
        nil, labels, // Using const labels
    )
    for _, name := range be.Cache.MetricNames() {
        m, err := be.Cache.Get(name)
        if err != nil {
            log.Errorf("Unable get metric %v for %v during GenericDescribe due to the error: %v", name, be.Name, err)
            continue
        }
        ch <- m.PromDescription(be.Name)
    }
}

/*
GenericCollect generates prom metric from cached metric values
*/
func GenericCollect(be *BaseExporter, ch chan<- prometheus.Metric) {
    log.Debugf("GenericCollect from %v %p \n", be.Name, be)
    for _, name := range be.Cache.MetricNames() {
        func(){
            m, err := be.Cache.Get(name)
            if err != nil {
                log.Warnf("Unable to get metric %v for %v during GenericCollect due to the error: %v", name, be.Name, err)
                return
            }
            metric, err := prometheus.NewConstMetric(
                m.PromDescription(be.Name),
                m.PromType(),
                m.Value,
            )
            if err != nil {
                log.Errorf("Unable to send metric %v for %v to the Prometheus client due to the error: %v", name, be.Name, err)
                return
            }
            ch <- metric
        }()
    }
}
