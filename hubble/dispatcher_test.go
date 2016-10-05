package hubble_test

import (
    . "github.com/aksentyev/hubble/hubble"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "testing"

    "fmt"
)

func TestDispatcher(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Dispather")
}

var _ = Describe("Dispather", func() {
    defer GinkgoRecover()

    Describe("Create new dispatcher instance", func() {
        It("should return new instance", func(){
            d := NewDispatcher(30)
            Expect(fmt.Sprintf("%T", d)).To(Equal("*hubble.Dispatcher"))
        })
    })
})
