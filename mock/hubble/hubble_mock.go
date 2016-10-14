package hubble_mock

import (
    "github.com/aksentyev/hubble/hubble"
    "fmt"
)

// Mock of BackendAdapter interface
type MockBackendAdapter struct {}
type MockKVBackendAdapter struct {}

func NewMockBackendAdapter() *MockBackendAdapter {
    mock := MockBackendAdapter{}
    return &mock
}

func NewMockKVBackendAdapter() *MockKVBackendAdapter {
    mock := MockKVBackendAdapter{}
    return &mock
}

func (m *MockBackendAdapter) GetAll() (list []*hubble.Service, err error) {
    for i := 0; i < 2; i++ {
        svc := hubble.DefaultService()
        svc.Name = fmt.Sprintf("testsvc%v", string(i))
        svc.Addresses = map[string]string{"srv1": fmt.Sprintf("10.0.0.1%v", string(i))}
        svc.Port = "9999"
        svc.Tags = []string{"test"}
        svc.ServiceParams.ModifyIndex = uint64(123456789) + uint64(i)
        list = append(list, svc)
    }
    return list, err
}

func (m *MockKVBackendAdapter) GetParams(_ string) (*hubble.ServiceParams, error) {
    return &hubble.ServiceParams{}, nil
}
