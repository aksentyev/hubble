package hubble_test

import (
    . "github.com/aksentyev/hubble/hubble"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "testing"
    "fmt"
    hubble_mock "../mock/hubble"


)

func TestHubble(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Hubble")
}

var _ = Describe("Hubble", func() {
    defer GinkgoRecover()

    cb := hubble_mock.NewMockBackendAdapter()
    kb := hubble_mock.NewMockKVBackendAdapter()


    Describe("Create new hubble instance", func() {
        It("should return new instance", func(){
            h := New(cb, kb, "testing")
            Expect(fmt.Sprintf("%T", h)).To(Equal("*hubble.Hubble"))
        })
    })

    Describe("Get Services", func() {
        h := New(cb, kb, "testing")
        filter1 := func(list []*Service) []*Service {
            return list[:len(list)-1]
        }
        filter2 := func(list []*Service) []*Service {
            return list
        }

        list, err := h.Services(filter1)
        all, err  := h.Services(filter2)

        It("should be 1p less than source services list", func(){
            Expect(len(list)).To(Equal(len(all)-1))
        })
        It("should not error", func() {
            Expect(err).NotTo(HaveOccurred())
        })
    })
})
