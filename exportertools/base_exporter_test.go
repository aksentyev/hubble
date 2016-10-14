package exportertools_test

import (
    . "github.com/aksentyev/hubble/exportertools"
    mock "../mock/exportertools"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "testing"

    "time"
    "fmt"
)

func TestBaseExporter(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "BaseExporter")
}

var _ = Describe("BaseExporter", func() {
    defer GinkgoRecover()

    var be *BaseExporter
    BeforeEach(func(){
        be = NewBaseExporter("test", 30, map[string]string{"key": "value"})
    })

    Describe("BaseExporter", func() {
        It("should return new instance", func(){
            Expect(fmt.Sprintf("%T", be)).To(Equal("*exportertools.BaseExporter"))
            Expect(be.Cache.TTL).To(Equal(30 * time.Second))
        })

        It("should add collector to list", func() {
            s1 := len(be.MetricCollectors)
            be.AddCollector(mock.NewMockCollector())
            s2 := len(be.MetricCollectors)

            Expect(s2).To(Equal(s1 + 1))
        })
    })
})
