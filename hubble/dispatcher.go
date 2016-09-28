package hubble

import (
    "time"
    "sync"
    hash "github.com/cnf/structhash"
    "encoding/hex"
    "github.com/aksentyev/hubble/exportertools"

    "github.com/prometheus/common/log"

    "os"
    "os/signal"
    "syscall"
)

/*
Dispatcher sends messages to ToRegister and ToUnregister channels
according to service list and its params. Dispatcher prevents multiple
registration for one service and unregistering for actual service.
Services and Exporter are stored in map where key is
*/
type Dispatcher struct {
    TTL          time.Duration
    Services     map[string]*ServiceAtomic
    Exporters    map[string]Exporter
    ToRegister   chan *ServiceAtomic
    ToUnregister chan map[string]*ServiceAtomic
    mtx          *sync.Mutex
}

/* Exporter interface is used to call Close() when Interrupt
or other quit signals are received
*/
type Exporter interface {
    exportertools.Exporter
}

/*
NewDispatcher receives TTL and callback function that should return []*ServiceAtomic.
Example:
    cb := func() (list []*hubble.ServiceAtomic, err error) {
        defer func() {
            if r := recover(); r != nil {
                err = errors.New(fmt.Sprintf("Unable to get services from consul: %v", r))
                list = []*hubble.ServiceAtomic{}
                log.Errorln(err)
            }
        }()

        for _, svc := range h.Services(filterCB){
            for _, el := range svc.MakeAtomic(nil) {
                list = append(list, el)
            }
        }
        return list, err
    }
    c = h.NewDispatcher(60, callback)

defer in example function is important. If callback was not return an error
Dispatcher receives empty services list. Re-register will be failed because prometheus golang client
cannot fully unregister Collector.
*/
func (h *Hubble) NewDispatcher(ttl int, cb func() ([]*ServiceAtomic, error)) *Dispatcher {
    d := Dispatcher{
        TTL:          time.Duration(ttl) * time.Second,
        Services:     map[string]*ServiceAtomic{},
        Exporters:    map[string]Exporter{},
        ToRegister:   make(chan *ServiceAtomic),
        ToUnregister: make(chan map[string]*ServiceAtomic),
        mtx:          &sync.Mutex{},
    }

    go d.run(cb)
    return &d
}

func (d *Dispatcher) run(f func() ([]*ServiceAtomic, error)) {
    ticker := time.NewTicker(d.TTL)

    // Gracefully exit on Interrupt
    go func() {
        c := make(chan os.Signal, 1)
        signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
        <-c
        log.Warnf("Termination signal received. Exiting.")
        for _, i := range d.Exporters {
            err := i.Close()
            if err != nil {
                log.Errorln(err)
            }
        }
        os.Exit(0)
    }()

    for _ = range ticker.C {
        func() {
            list, err := f()
            if err != nil {
                return
            }

            actual := map[string]*ServiceAtomic{}
            for _, svc := range list {
                h := hex.EncodeToString(hash.Md5(svc, 1))
                actual[h] = svc
            }

            toBeRemoved := map[string]*ServiceAtomic{}
            for k, v := range d.Services {
                toBeRemoved[k] = v
            }

            for h, svc := range actual {
                if _, ok := d.Services[h]; !ok {
                    go func(s *ServiceAtomic) { d.ToRegister <- s }(svc)
                }
                delete(toBeRemoved, h)
            }
            for h, svc := range toBeRemoved {
                d.ToUnregister <- map[string]*ServiceAtomic{h: svc}
            }
        }()
    }
}

// Register service and exporter pair in Dispatcher
func (d *Dispatcher) Register(item *ServiceAtomic, exporter Exporter) {
    d.mtx.Lock()
    defer d.mtx.Unlock()
    h := hex.EncodeToString(hash.Md5(item, 1))
    d.Services[h] = item
    d.Exporters[h] = exporter
}

// Unregister service in Dispatcher
func (d *Dispatcher) Unregister(item *ServiceAtomic) {
    d.mtx.Lock()
    defer d.mtx.Unlock()

    h := hex.EncodeToString(hash.Md5(item, 1))
    delete(d.Services, h)
    delete(d.Exporters, h)
}

// Unregister service in Dispatcher with only hash identificator provided
func (d *Dispatcher) UnregisterWithHash(h string) {
    d.mtx.Lock()
    defer d.mtx.Unlock()

    delete(d.Services, h)
    delete(d.Exporters, h)
}
