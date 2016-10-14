package exportertools_test

import (
    . "github.com/aksentyev/hubble/exportertools"
    mock "../mock/exportertools"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "testing"
)

func TestCollector(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Collector")
}

var _ = Describe("Collector", func() {
    defer GinkgoRecover()

    var be *BaseExporter
    BeforeEach(func(){
        be = NewBaseExporter("test", 30, map[string]string{"key": "value"})
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
