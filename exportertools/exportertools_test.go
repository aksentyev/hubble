package exportertools_test

import (
    . "github.com/aksentyev/hubble/exportertools"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "testing"
    "time"

    "fmt"

    mock "../mock/exportertools"
)

func TestBaseExporter(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Exportertools")
}

var _ = Describe("Exportertools", func() {
    defer GinkgoRecover()

    var be *BaseExporter

    BeforeSuite(func(){
        be = NewBaseExporter("test", 30, map[string]string{"key": "value"})
    })

    Context("BaseExporter", func() {
        Describe("Create new BaseExporter instance", func() {
            It("should return new instance", func(){
                Expect(fmt.Sprintf("%T", be)).To(Equal("*exportertools.BaseExporter"))
                Expect(be.Cache.TTL).To(Equal(30 * time.Second))
            })
        })
    })

    Describe("Collector", func() {
        var metric *Metric
        It("should collect metrics",func(){
            m := mock.NewMockCollector()
            be.AddCollector(m)
            metrics, _ := be.MetricCollectors[0].Collect()
            metric = metrics[0]

            Expect(len(be.MetricCollectors)).To(Equal(1))
            Expect(metric.Name).To(Equal("test_collector_metric"))
        })
    })
})
