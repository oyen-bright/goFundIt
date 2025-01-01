// Code generated by mockery v2.50.0. DO NOT EDIT.

package encryption

import (
	encryption "github.com/oyen-bright/goFundIt/pkg/encryption"
	mock "github.com/stretchr/testify/mock"
)

// MockEncryptor is an autogenerated mock type for the Encryptor type
type MockEncryptor struct {
	mock.Mock
}

type MockEncryptor_Expecter struct {
	mock *mock.Mock
}

func (_m *MockEncryptor) EXPECT() *MockEncryptor_Expecter {
	return &MockEncryptor_Expecter{mock: &_m.Mock}
}

// Decrypt provides a mock function with given fields: data
func (_m *MockEncryptor) Decrypt(data encryption.Data) (string, error) {
	ret := _m.Called(data)

	if len(ret) == 0 {
		panic("no return value specified for Decrypt")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(encryption.Data) (string, error)); ok {
		return rf(data)
	}
	if rf, ok := ret.Get(0).(func(encryption.Data) string); ok {
		r0 = rf(data)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(encryption.Data) error); ok {
		r1 = rf(data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockEncryptor_Decrypt_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Decrypt'
type MockEncryptor_Decrypt_Call struct {
	*mock.Call
}

// Decrypt is a helper method to define mock.On call
//   - data encryption.Data
func (_e *MockEncryptor_Expecter) Decrypt(data interface{}) *MockEncryptor_Decrypt_Call {
	return &MockEncryptor_Decrypt_Call{Call: _e.mock.On("Decrypt", data)}
}

func (_c *MockEncryptor_Decrypt_Call) Run(run func(data encryption.Data)) *MockEncryptor_Decrypt_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(encryption.Data))
	})
	return _c
}

func (_c *MockEncryptor_Decrypt_Call) Return(_a0 string, _a1 error) *MockEncryptor_Decrypt_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockEncryptor_Decrypt_Call) RunAndReturn(run func(encryption.Data) (string, error)) *MockEncryptor_Decrypt_Call {
	_c.Call.Return(run)
	return _c
}

// DecryptStruct provides a mock function with given fields: data, key
func (_m *MockEncryptor) DecryptStruct(data interface{}, key string) (interface{}, error) {
	ret := _m.Called(data, key)

	if len(ret) == 0 {
		panic("no return value specified for DecryptStruct")
	}

	var r0 interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(interface{}, string) (interface{}, error)); ok {
		return rf(data, key)
	}
	if rf, ok := ret.Get(0).(func(interface{}, string) interface{}); ok {
		r0 = rf(data, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(interface{}, string) error); ok {
		r1 = rf(data, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockEncryptor_DecryptStruct_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DecryptStruct'
type MockEncryptor_DecryptStruct_Call struct {
	*mock.Call
}

// DecryptStruct is a helper method to define mock.On call
//   - data interface{}
//   - key string
func (_e *MockEncryptor_Expecter) DecryptStruct(data interface{}, key interface{}) *MockEncryptor_DecryptStruct_Call {
	return &MockEncryptor_DecryptStruct_Call{Call: _e.mock.On("DecryptStruct", data, key)}
}

func (_c *MockEncryptor_DecryptStruct_Call) Run(run func(data interface{}, key string)) *MockEncryptor_DecryptStruct_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}), args[1].(string))
	})
	return _c
}

func (_c *MockEncryptor_DecryptStruct_Call) Return(_a0 interface{}, _a1 error) *MockEncryptor_DecryptStruct_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockEncryptor_DecryptStruct_Call) RunAndReturn(run func(interface{}, string) (interface{}, error)) *MockEncryptor_DecryptStruct_Call {
	_c.Call.Return(run)
	return _c
}

// Encrypt provides a mock function with given fields: data
func (_m *MockEncryptor) Encrypt(data encryption.Data) (string, error) {
	ret := _m.Called(data)

	if len(ret) == 0 {
		panic("no return value specified for Encrypt")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(encryption.Data) (string, error)); ok {
		return rf(data)
	}
	if rf, ok := ret.Get(0).(func(encryption.Data) string); ok {
		r0 = rf(data)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(encryption.Data) error); ok {
		r1 = rf(data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockEncryptor_Encrypt_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Encrypt'
type MockEncryptor_Encrypt_Call struct {
	*mock.Call
}

// Encrypt is a helper method to define mock.On call
//   - data encryption.Data
func (_e *MockEncryptor_Expecter) Encrypt(data interface{}) *MockEncryptor_Encrypt_Call {
	return &MockEncryptor_Encrypt_Call{Call: _e.mock.On("Encrypt", data)}
}

func (_c *MockEncryptor_Encrypt_Call) Run(run func(data encryption.Data)) *MockEncryptor_Encrypt_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(encryption.Data))
	})
	return _c
}

func (_c *MockEncryptor_Encrypt_Call) Return(_a0 string, _a1 error) *MockEncryptor_Encrypt_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockEncryptor_Encrypt_Call) RunAndReturn(run func(encryption.Data) (string, error)) *MockEncryptor_Encrypt_Call {
	_c.Call.Return(run)
	return _c
}

// EncryptStruct provides a mock function with given fields: data, key
func (_m *MockEncryptor) EncryptStruct(data interface{}, key string) (interface{}, error) {
	ret := _m.Called(data, key)

	if len(ret) == 0 {
		panic("no return value specified for EncryptStruct")
	}

	var r0 interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(interface{}, string) (interface{}, error)); ok {
		return rf(data, key)
	}
	if rf, ok := ret.Get(0).(func(interface{}, string) interface{}); ok {
		r0 = rf(data, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(interface{}, string) error); ok {
		r1 = rf(data, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockEncryptor_EncryptStruct_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EncryptStruct'
type MockEncryptor_EncryptStruct_Call struct {
	*mock.Call
}

// EncryptStruct is a helper method to define mock.On call
//   - data interface{}
//   - key string
func (_e *MockEncryptor_Expecter) EncryptStruct(data interface{}, key interface{}) *MockEncryptor_EncryptStruct_Call {
	return &MockEncryptor_EncryptStruct_Call{Call: _e.mock.On("EncryptStruct", data, key)}
}

func (_c *MockEncryptor_EncryptStruct_Call) Run(run func(data interface{}, key string)) *MockEncryptor_EncryptStruct_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}), args[1].(string))
	})
	return _c
}

func (_c *MockEncryptor_EncryptStruct_Call) Return(_a0 interface{}, _a1 error) *MockEncryptor_EncryptStruct_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockEncryptor_EncryptStruct_Call) RunAndReturn(run func(interface{}, string) (interface{}, error)) *MockEncryptor_EncryptStruct_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockEncryptor creates a new instance of MockEncryptor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockEncryptor(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockEncryptor {
	mock := &MockEncryptor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
