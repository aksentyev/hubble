package main
//
// import (
//     "github.com/aksentyev/hubble/hubble"
//     "github.com/aksentyev/hubble/consul"
//     "fmt"
//     "encoding/json"
// )
//
// func main() {
//     params := &consul.Config{
//         Address: "consul.service.consul:8500",
//         Datacenter: "staging",
//     }
//     client, _ := consul.New(params)
//
//     kv := consul.NewKV(client)
//
//     h := hubble.New(client, kv)
//     list := h.Services()
//
//     for _, svc := range list{
//         b, _ := json.MarshalIndent(svc, "", "  ")
//         fmt.Printf("Service:\n%+v\n", string(b))
//
//          l := svc.MakeAtomic(nil)
//         b, _ = json.MarshalIndent(l, "", "  ")
//         fmt.Printf("ServiceAtomic:\n%+v\n", string(b))
//     }
//
// }
