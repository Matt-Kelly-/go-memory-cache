// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package store

import mock "github.com/stretchr/testify/mock"

// MockStore is an autogenerated mock type for the Store type
type MockStore struct {
	mock.Mock
}

// Delete provides a mock function with given fields: key
func (_m *MockStore) Delete(key string) {
	_m.Called(key)
}

// Get provides a mock function with given fields: key
func (_m *MockStore) Get(key string) (string, bool) {
	ret := _m.Called(key)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// Has provides a mock function with given fields: key
func (_m *MockStore) Has(key string) bool {
	ret := _m.Called(key)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Put provides a mock function with given fields: key, value
func (_m *MockStore) Put(key string, value string) {
	_m.Called(key, value)
}
