package hubble

// Service main parameters
type Service struct {
    ModifyIndex uint64
    Name        string
    Addresses   map[string]string
    Port        string
    Tags        []string
    *ServiceParams
}

// ServiceParams struct to keep other service check or metric parameters.
type ServiceParams struct {
    ModifyIndex     uint64
    Notifiable      bool              `json:"notifiable"` // useful for Alertmanager
    ExtraLabels     map[string]string `json:"extra_labels"`
    ExporterOptions map[string]string `json:"exporter_options"`
}

/*
 ServiceAtomic is the main useful unit for an expoter.
Actually it is Service splitted by host.
*/
type ServiceAtomic struct {
    Name            string
    Address         string
    Port            string
    Tags           []string
    ExtraLabels     map[string]string
    ExporterOptions map[string]string
}
