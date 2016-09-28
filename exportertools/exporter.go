package exportertools

import (
    "github.com/prometheus/client_golang/prometheus"

    "errors"
    "fmt"
)

type Exporter interface {
    Setup() error
    Close() error

    Process()

    prometheus.Collector
}

func Register(exporter Exporter) (err error) {
    err = exporter.Setup()
    if err != nil {
        return err
    }

    err = prometheus.Register(exporter)
    if err != nil {
        return err
    }
    go exporter.Process()
    return nil
}

func Unregister(exporter Exporter) error {
    ok := prometheus.Unregister(exporter)
    if !ok {
        return errors.New(fmt.Sprintf("Unregister() failed for %v", exporter))
    }
    return nil
}
