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
    *sync.RWMutex
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

    d := hubble.NewDispatcher(60)
    d.Run(callback)

defer in example function is important. If callback was not return an error
Dispatcher receives empty services list. Re-register will be failed because prometheus golang client
cannot fully unregister Collector.
*/
func NewDispatcher(ttl int) *Dispatcher {
    d := Dispatcher{
        TTL:          time.Duration(ttl) * time.Second,
        Services:     map[string]*ServiceAtomic{},
        Exporters:    map[string]Exporter{},
        ToRegister:   make(chan *ServiceAtomic, 20),
        ToUnregister: make(chan map[string]*ServiceAtomic, 20),
        RWMutex:      &sync.RWMutex{},
    }

    return &d
}

// Run Dispatcher
func (d *Dispatcher) Run(f func() ([]*ServiceAtomic, error)) {
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
        d.process(f)
    }
}

func (d *Dispatcher) process( f func() ([]*ServiceAtomic, error) ) {
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
    d.RLock()
    for k, v := range d.Services {
        toBeRemoved[k] = v
    }
    d.RUnlock()

    for h, svc := range actual {
        d.RLock()
        _, ok := d.Services[h]
        d.RUnlock()

        if !ok {
            d.ToRegister <- svc
        }
        delete(toBeRemoved, h)
    }

    for h, svc := range toBeRemoved {
        d.ToUnregister <- map[string]*ServiceAtomic{h: svc}
    }
}

// Register service and exporter pair in Dispatcher
func (d *Dispatcher) Register(item *ServiceAtomic, exporter Exporter) {
    h := hex.EncodeToString(hash.Md5(item, 1))

    d.Lock()
    defer d.Unlock()
    d.Services[h] = item
    d.Exporters[h] = exporter
}

// Unregister service in Dispatcher
func (d *Dispatcher) Unregister(item *ServiceAtomic) {
    h := hex.EncodeToString(hash.Md5(item, 1))

    d.Lock()
    defer d.Unlock()
    delete(d.Services, h)
    delete(d.Exporters, h)
}

// Unregister service in Dispatcher with only hash identificator provided
func (d *Dispatcher) UnregisterWithHash(h string) {
    d.Lock()
    defer d.Unlock()

    delete(d.Services, h)
    delete(d.Exporters, h)
}
