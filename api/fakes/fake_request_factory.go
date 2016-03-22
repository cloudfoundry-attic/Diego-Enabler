package fakes

import (
	"net/http"
	"sync"

	"github.com/cloudfoundry-incubator/diego-enabler/api"
)

type FakeRequestFactory struct {
	Stub        func(api.Filter, map[string]interface{}) (*http.Request, error)
	mutex       sync.RWMutex
	argsForCall []struct {
		arg1 api.Filter
		arg2 map[string]interface{}
	}
	returns struct {
		result1 *http.Request
		result2 error
	}
}

func (fake *FakeRequestFactory) Spy(arg1 api.Filter, arg2 map[string]interface{}) (*http.Request, error) {
	fake.mutex.Lock()
	fake.argsForCall = append(fake.argsForCall, struct {
		arg1 api.Filter
		arg2 map[string]interface{}
	}{arg1, arg2})
	fake.mutex.Unlock()
	if fake.Stub != nil {
		return fake.Stub(arg1, arg2)
	} else {
		return fake.returns.result1, fake.returns.result2
	}
}

func (fake *FakeRequestFactory) CallCount() int {
	fake.mutex.RLock()
	defer fake.mutex.RUnlock()
	return len(fake.argsForCall)
}

func (fake *FakeRequestFactory) ArgsForCall(i int) (api.Filter, map[string]interface{}) {
	fake.mutex.RLock()
	defer fake.mutex.RUnlock()
	return fake.argsForCall[i].arg1, fake.argsForCall[i].arg2
}

func (fake *FakeRequestFactory) Returns(result1 *http.Request, result2 error) {
	fake.Stub = nil
	fake.returns = struct {
		result1 *http.Request
		result2 error
	}{result1, result2}
}

var _ api.RequestFactory = new(FakeRequestFactory).Spy
