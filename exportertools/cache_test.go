package exportertools_test

import (
    . "github.com/aksentyev/hubble/exportertools"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "testing"

    "time"
    "fmt"
)

func TestCache(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Cache")
}

var _ = Describe("Cache", func() {
    defer GinkgoRecover()

    var cache *Cache
    BeforeEach(func() {
        cache = NewCache(30)
    })

    Describe("Cache instance", func() {
        It("should return new instance", func(){
            Expect(fmt.Sprintf("%T", cache)).To(Equal("*exportertools.Cache"))
            Expect(cache.TTL).To(Equal(30 * time.Second))
        })

        It("should store and return metrics", func(){
            m := &Metric{}
            m.Name = "Test"
            cache.Set(m)
            Expect(cache.Get("Test")).To(Equal(m))
        })

        It("should metrics names list", func(){
            m := &Metric{}
            m.Name = "Test"
            cache.Set(m)
            Expect(cache.MetricNames()).To(Equal([]string{"Test"}))
        })
    })
})
