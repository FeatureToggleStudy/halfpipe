// Code generated by counterfeiter. DO NOT EDIT.
package dbfakes

import (
	"sync"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
	"github.com/concourse/atc/creds"
	"github.com/concourse/atc/db"
)

type FakeResourceConfigCheckSessionFactory struct {
	FindOrCreateResourceConfigCheckSessionStub        func(logger lager.Logger, resourceType string, source atc.Source, resourceTypes creds.VersionedResourceTypes, expiries db.ContainerOwnerExpiries) (db.ResourceConfigCheckSession, error)
	findOrCreateResourceConfigCheckSessionMutex       sync.RWMutex
	findOrCreateResourceConfigCheckSessionArgsForCall []struct {
		logger        lager.Logger
		resourceType  string
		source        atc.Source
		resourceTypes creds.VersionedResourceTypes
		expiries      db.ContainerOwnerExpiries
	}
	findOrCreateResourceConfigCheckSessionReturns struct {
		result1 db.ResourceConfigCheckSession
		result2 error
	}
	findOrCreateResourceConfigCheckSessionReturnsOnCall map[int]struct {
		result1 db.ResourceConfigCheckSession
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeResourceConfigCheckSessionFactory) FindOrCreateResourceConfigCheckSession(logger lager.Logger, resourceType string, source atc.Source, resourceTypes creds.VersionedResourceTypes, expiries db.ContainerOwnerExpiries) (db.ResourceConfigCheckSession, error) {
	fake.findOrCreateResourceConfigCheckSessionMutex.Lock()
	ret, specificReturn := fake.findOrCreateResourceConfigCheckSessionReturnsOnCall[len(fake.findOrCreateResourceConfigCheckSessionArgsForCall)]
	fake.findOrCreateResourceConfigCheckSessionArgsForCall = append(fake.findOrCreateResourceConfigCheckSessionArgsForCall, struct {
		logger        lager.Logger
		resourceType  string
		source        atc.Source
		resourceTypes creds.VersionedResourceTypes
		expiries      db.ContainerOwnerExpiries
	}{logger, resourceType, source, resourceTypes, expiries})
	fake.recordInvocation("FindOrCreateResourceConfigCheckSession", []interface{}{logger, resourceType, source, resourceTypes, expiries})
	fake.findOrCreateResourceConfigCheckSessionMutex.Unlock()
	if fake.FindOrCreateResourceConfigCheckSessionStub != nil {
		return fake.FindOrCreateResourceConfigCheckSessionStub(logger, resourceType, source, resourceTypes, expiries)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.findOrCreateResourceConfigCheckSessionReturns.result1, fake.findOrCreateResourceConfigCheckSessionReturns.result2
}

func (fake *FakeResourceConfigCheckSessionFactory) FindOrCreateResourceConfigCheckSessionCallCount() int {
	fake.findOrCreateResourceConfigCheckSessionMutex.RLock()
	defer fake.findOrCreateResourceConfigCheckSessionMutex.RUnlock()
	return len(fake.findOrCreateResourceConfigCheckSessionArgsForCall)
}

func (fake *FakeResourceConfigCheckSessionFactory) FindOrCreateResourceConfigCheckSessionArgsForCall(i int) (lager.Logger, string, atc.Source, creds.VersionedResourceTypes, db.ContainerOwnerExpiries) {
	fake.findOrCreateResourceConfigCheckSessionMutex.RLock()
	defer fake.findOrCreateResourceConfigCheckSessionMutex.RUnlock()
	return fake.findOrCreateResourceConfigCheckSessionArgsForCall[i].logger, fake.findOrCreateResourceConfigCheckSessionArgsForCall[i].resourceType, fake.findOrCreateResourceConfigCheckSessionArgsForCall[i].source, fake.findOrCreateResourceConfigCheckSessionArgsForCall[i].resourceTypes, fake.findOrCreateResourceConfigCheckSessionArgsForCall[i].expiries
}

func (fake *FakeResourceConfigCheckSessionFactory) FindOrCreateResourceConfigCheckSessionReturns(result1 db.ResourceConfigCheckSession, result2 error) {
	fake.FindOrCreateResourceConfigCheckSessionStub = nil
	fake.findOrCreateResourceConfigCheckSessionReturns = struct {
		result1 db.ResourceConfigCheckSession
		result2 error
	}{result1, result2}
}

func (fake *FakeResourceConfigCheckSessionFactory) FindOrCreateResourceConfigCheckSessionReturnsOnCall(i int, result1 db.ResourceConfigCheckSession, result2 error) {
	fake.FindOrCreateResourceConfigCheckSessionStub = nil
	if fake.findOrCreateResourceConfigCheckSessionReturnsOnCall == nil {
		fake.findOrCreateResourceConfigCheckSessionReturnsOnCall = make(map[int]struct {
			result1 db.ResourceConfigCheckSession
			result2 error
		})
	}
	fake.findOrCreateResourceConfigCheckSessionReturnsOnCall[i] = struct {
		result1 db.ResourceConfigCheckSession
		result2 error
	}{result1, result2}
}

func (fake *FakeResourceConfigCheckSessionFactory) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.findOrCreateResourceConfigCheckSessionMutex.RLock()
	defer fake.findOrCreateResourceConfigCheckSessionMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeResourceConfigCheckSessionFactory) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ db.ResourceConfigCheckSessionFactory = new(FakeResourceConfigCheckSessionFactory)