package hubble

import (
    "github.com/prometheus/common/log"
    "reflect"
    "sync"
    "fmt"
    "strings"
    "errors"
)

// Hubble stores service and kv backends adapters
type Hubble struct {
    SvcBackend    BackendAdapter
    ParamsBackend KVBackendAdapter
    exporterName string
}

// BackendAdapter implements GetAll() method
type BackendAdapter interface {
    GetAll() ([]*Service, error)
}

// KVBackendAdapter implements GetParams() method should fill ServiceParams part of Service with data
type KVBackendAdapter interface {
    GetParams(key string) (*ServiceParams, error)
}

func New(sb BackendAdapter, pb KVBackendAdapter, exporterName string) *Hubble {
    h := &Hubble{sb, pb, exporterName}
    return h
}

/*
Services() recever a filter callback to filter services
will not be used in exporter. Then returns []*Service.

example:
    filterCB := func(list []*hubble.Service) []*hubble.Service {
        var servicesForMonitoring []*hubble.Service
        for _, svc := range list {
            if util.IncludesStr(svc.Tags, "goro") {
                servicesForMonitoring = append(servicesForMonitoring, svc)
            }
        }
        return servicesForMonitoring
    }
*/
func (h *Hubble) Services(filter func(list []*Service) []*Service) ([]*Service, error) {
    var servicesForMonitoring []*Service

    log.Debugf("Getting service from backend")
    allServices, err := h.getServices()
    if err == nil {
        log.Debugf("Service list %v\n", names(allServices))
        servicesForMonitoring = filter(allServices)

        log.Debugf("Getting service params from backend")
        err = h.getParams(servicesForMonitoring)

        if err == nil {
            log.Debugf("Services, with defined params were found: %v\n", notDefaultParams(servicesForMonitoring))
        }
    }

    return servicesForMonitoring, err
}

func (h *Hubble) getServices() ([]*Service, error) {
    return h.SvcBackend.GetAll()
}

func (h *Hubble) getParams(list []*Service) (err error) {
    var wg sync.WaitGroup
    threads := make(chan bool, 10)
    errCh := make(chan error)

    wg.Add(len(list))

    go func() {
        defer close(threads)
        defer close(errCh)

        for _, svc := range list {
            threads<- true
            go func(svc *Service) {
                defer func() { <-threads }()
                defer wg.Done()

                params, err := func(svc *Service) (*ServiceParams, error) {
                    key := fmt.Sprintf("monitoring/%v/%v", strings.ToLower(svc.Name), strings.ToLower(h.exporterName))
                    params, err := h.ParamsBackend.GetParams(key)
                    if err != nil {
                        err = errors.New(fmt.Sprintf("Get params for svc %v was failed: %v", svc.Name, err))
                        return &ServiceParams{}, err
                    }
                    return params, nil
                }(svc)

                if err != nil {
                    log.Errorln(err)
                    errCh <- err
                    return
                }
                svc.ServiceParams = params
            }(svc)
        }
        wg.Wait()
    }()

    err = <- errCh
    return err
}

func names(services []*Service) (list []string) {
    for _, s := range services {
        list = append(list, s.Name)
    }
    return list
}

func notDefaultParams(services []*Service) (list []string) {
    defaultService := DefaultService()
    for _, s := range services {
        p1 := s.ServiceParams
        p2 := defaultService.ServiceParams

        modified := p1.Notifiable != p2.Notifiable ||
                    reflect.DeepEqual(p1.ExtraLabels, p2.ExtraLabels) ||
                    reflect.DeepEqual(p1.ExporterOptions, p2.ExporterOptions)

        if modified {
            list = append(list, s.Name)
        }
    }
    return list
}
