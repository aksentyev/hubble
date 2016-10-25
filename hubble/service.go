package hubble

import (
    "strconv"
)

// DefaultService return empty Service instance
func DefaultService() *Service {
    s := &Service{
        Name: "",
        Addresses: map[string]string{},
        Port: "",
        Tags: []string{},
        ServiceParams: &ServiceParams{
            Notifiable: false,
            ExtraLabels: map[string]string{},
            ExporterOptions: map[string]string{},
        },
    }

    return s
}

// MakeAtomic splits Service into list of ServiceAtomic
func (s *Service) MakeAtomic(extraLabels map[string]string) []*ServiceAtomic {
    var list []*ServiceAtomic
    for host := range s.Addresses {
        atomic := &ServiceAtomic{
            Name:        s.Name,
            Address:     s.Addresses[host],
            Port:        s.Port,
            Tags:        s.Tags,
            ExtraLabels: map[string]string{
                "service": s.Name,
                "host": host,
                "notifiable": strconv.FormatBool(s.Notifiable),
                "modify_index": strconv.FormatUint(s.ModifyIndex + s.ServiceParams.ModifyIndex, 10),
            },
        }

        for key, val := range extraLabels {
            atomic.ExtraLabels[key] = val
        }

        for key, val := range s.ExtraLabels {
            atomic.ExtraLabels[key] = val
        }

        atomic.ExporterOptions = s.ExporterOptions
        list = append(list, atomic)
    }
    return list
}
