package tester

import (
	cabby "github.com/pladdy/cabby2"
)

/* DataStore */

// DataStore for tests
type DataStore struct {
	APIRootServiceFn    func() APIRootService
	CollectionServiceFn func() CollectionService
	DiscoveryServiceFn  func() DiscoveryService
	UserServiceFn       func() UserService
}

// NewDataStore structure
func NewDataStore() *DataStore {
	return &DataStore{}
}

// APIRootService mock
func (s DataStore) APIRootService() cabby.APIRootService {
	return s.APIRootServiceFn()
}

// Close mock
func (s DataStore) Close() {
	return
}

// CollectionService mock
func (s DataStore) CollectionService() cabby.CollectionService {
	return s.CollectionServiceFn()
}

// DiscoveryService mock
func (s DataStore) DiscoveryService() cabby.DiscoveryService {
	return s.DiscoveryServiceFn()
}

// Open mock
func (s DataStore) Open() error {
	return nil
}

// UserService mock
func (s DataStore) UserService() cabby.UserService {
	return s.UserServiceFn()
}

/* services */

// APIRootService is a mock implementation
type APIRootService struct {
	APIRootFn  func(path string) (cabby.APIRoot, error)
	APIRootsFn func() ([]cabby.APIRoot, error)
}

// APIRoot is a mock implementation
func (s APIRootService) APIRoot(path string) (cabby.APIRoot, error) {
	return s.APIRootFn(path)
}

// APIRoots is a mock implementation
func (s APIRootService) APIRoots() ([]cabby.APIRoot, error) {
	return s.APIRootsFn()
}

// CollectionService is a mock implementation
type CollectionService struct {
	CollectionFn           func(user, collectionID, apiRootPath string) (cabby.Collection, error)
	CollectionsFn          func(user, apiRootPath string) (cabby.Collections, error)
	CollectionsInAPIRootFn func(apiRootPath string) (cabby.CollectionsInAPIRoot, error)
}

// Collection is a mock implementation
func (s CollectionService) Collection(user, collectionID, apiRootPath string) (cabby.Collection, error) {
	return s.CollectionFn(user, collectionID, apiRootPath)
}

// Collections is a mock implementation
func (s CollectionService) Collections(user, apiRootPath string) (cabby.Collections, error) {
	return s.CollectionsFn(user, apiRootPath)
}

// CollectionsInAPIRoot is a mock implementation
func (s CollectionService) CollectionsInAPIRoot(apiRootPath string) (cabby.CollectionsInAPIRoot, error) {
	return s.CollectionsInAPIRootFn(apiRootPath)
}

// DiscoveryService is a mock implementation
type DiscoveryService struct {
	DiscoveryFn func() (cabby.Discovery, error)
}

// Discovery is a mock implementation
func (s DiscoveryService) Discovery() (cabby.Discovery, error) {
	return s.DiscoveryFn()
}

// UserService is a mock implementation
type UserService struct {
	UserFn   func(user, password string) (cabby.User, error)
	ExistsFn func(cabby.User) bool
}

// User is a mock implementation
func (s UserService) User(user, password string) (cabby.User, error) {
	return s.UserFn(user, password)
}

// Exists is a mock implementation
func (s UserService) Exists(u cabby.User) bool {
	return s.ExistsFn(u)
}