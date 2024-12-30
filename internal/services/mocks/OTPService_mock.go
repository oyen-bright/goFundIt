// Code generated by mockery v2.50.0. DO NOT EDIT.

package interfaces

import (
	models "github.com/oyen-bright/goFundIt/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// MockOTPService is an autogenerated mock type for the OTPService type
type MockOTPService struct {
	mock.Mock
}

type MockOTPService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockOTPService) EXPECT() *MockOTPService_Expecter {
	return &MockOTPService_Expecter{mock: &_m.Mock}
}

// RequestOTP provides a mock function with given fields: email, name
func (_m *MockOTPService) RequestOTP(email string, name string) (models.Otp, error) {
	ret := _m.Called(email, name)

	if len(ret) == 0 {
		panic("no return value specified for RequestOTP")
	}

	var r0 models.Otp
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (models.Otp, error)); ok {
		return rf(email, name)
	}
	if rf, ok := ret.Get(0).(func(string, string) models.Otp); ok {
		r0 = rf(email, name)
	} else {
		r0 = ret.Get(0).(models.Otp)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(email, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOTPService_RequestOTP_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RequestOTP'
type MockOTPService_RequestOTP_Call struct {
	*mock.Call
}

// RequestOTP is a helper method to define mock.On call
//   - email string
//   - name string
func (_e *MockOTPService_Expecter) RequestOTP(email interface{}, name interface{}) *MockOTPService_RequestOTP_Call {
	return &MockOTPService_RequestOTP_Call{Call: _e.mock.On("RequestOTP", email, name)}
}

func (_c *MockOTPService_RequestOTP_Call) Run(run func(email string, name string)) *MockOTPService_RequestOTP_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *MockOTPService_RequestOTP_Call) Return(_a0 models.Otp, _a1 error) *MockOTPService_RequestOTP_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOTPService_RequestOTP_Call) RunAndReturn(run func(string, string) (models.Otp, error)) *MockOTPService_RequestOTP_Call {
	_c.Call.Return(run)
	return _c
}

// VerifyOTP provides a mock function with given fields: email, otp, requestId
func (_m *MockOTPService) VerifyOTP(email string, otp string, requestId string) (models.Otp, error) {
	ret := _m.Called(email, otp, requestId)

	if len(ret) == 0 {
		panic("no return value specified for VerifyOTP")
	}

	var r0 models.Otp
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, string) (models.Otp, error)); ok {
		return rf(email, otp, requestId)
	}
	if rf, ok := ret.Get(0).(func(string, string, string) models.Otp); ok {
		r0 = rf(email, otp, requestId)
	} else {
		r0 = ret.Get(0).(models.Otp)
	}

	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(email, otp, requestId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOTPService_VerifyOTP_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'VerifyOTP'
type MockOTPService_VerifyOTP_Call struct {
	*mock.Call
}

// VerifyOTP is a helper method to define mock.On call
//   - email string
//   - otp string
//   - requestId string
func (_e *MockOTPService_Expecter) VerifyOTP(email interface{}, otp interface{}, requestId interface{}) *MockOTPService_VerifyOTP_Call {
	return &MockOTPService_VerifyOTP_Call{Call: _e.mock.On("VerifyOTP", email, otp, requestId)}
}

func (_c *MockOTPService_VerifyOTP_Call) Run(run func(email string, otp string, requestId string)) *MockOTPService_VerifyOTP_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockOTPService_VerifyOTP_Call) Return(_a0 models.Otp, _a1 error) *MockOTPService_VerifyOTP_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOTPService_VerifyOTP_Call) RunAndReturn(run func(string, string, string) (models.Otp, error)) *MockOTPService_VerifyOTP_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockOTPService creates a new instance of MockOTPService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockOTPService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockOTPService {
	mock := &MockOTPService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
