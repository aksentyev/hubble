package hubble_test

import (
    . "github.com/aksentyev/hubble/hubble"
    hubble_mock "../mock/hubble"

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

    var h *Hubble

    BeforeSuite(func() {
        cb := hubble_mock.NewMockBackendAdapter()
        kb := hubble_mock.NewMockKVBackendAdapter()
        h = New(cb, kb, "testing")

    })


    Describe("Create new dispatcher instance", func() {
        It("should return new instance", func(){
            d := NewDispatcher(30)
            Expect(fmt.Sprintf("%T", d)).To(Equal("*hubble.Dispatcher"))
        })
    })

    Context("Dispatch services", func(){
        d := NewDispatcher(1)

        Describe("Services channels", func(){
            It("should receive services to be registered from channel", func(){
                filterCB := func(list []*Service) []*Service {
                    return list
                }

                cb := func() (list []*ServiceAtomic, err error) {
                    services, err := h.Services(filterCB)
                    if err != nil {
                        return list, err
                    }
                    for _, svc := range services {
                        for _, el := range svc.MakeAtomic(nil) {
                            list = append(list, el)
                        }
                    }
                    return list, err
                }
                d.Process(cb)

                Eventually(d.ToRegister).Should(Receive())
            })

            It("should receive services to be unregistered from channel", func(){
                filterCB := func(list []*Service) []*Service {
                    return []*Service{}
                }

                cb := func() (list []*ServiceAtomic, err error) {
                    services, err := h.Services(filterCB)
                    if err != nil {
                        return list, err
                    }
                    for _, svc := range services {
                        for _, el := range svc.MakeAtomic(nil) {
                            list = append(list, el)
                        }
                    }
                    return list, err
                }

                d.Register(&ServiceAtomic{}, nil)
                d.Process(cb)

                Eventually(d.ToUnregister).Should(Receive())
            })
        })
    })
})
