package hubble_mock

import "github.com/aksentyev/hubble/hubble"

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
        list = append(list, svc)
    }
    return list, err
}

func (m *MockKVBackendAdapter) GetParams(_ string) (*hubble.ServiceParams, error) {
    return &hubble.ServiceParams{}, nil
}
