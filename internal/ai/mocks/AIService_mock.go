// Code generated by mockery v2.50.0. DO NOT EDIT.

package interfaces

import (
	models "github.com/oyen-bright/goFundIt/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// MockAIService is an autogenerated mock type for the AIService type
type MockAIService struct {
	mock.Mock
}

type MockAIService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAIService) EXPECT() *MockAIService_Expecter {
	return &MockAIService_Expecter{mock: &_m.Mock}
}

// GenerateActivitySuggestions provides a mock function with given fields: campaignDescription
func (_m *MockAIService) GenerateActivitySuggestions(campaignDescription string) ([]models.ActivitySuggestion, error) {
	ret := _m.Called(campaignDescription)

	if len(ret) == 0 {
		panic("no return value specified for GenerateActivitySuggestions")
	}

	var r0 []models.ActivitySuggestion
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]models.ActivitySuggestion, error)); ok {
		return rf(campaignDescription)
	}
	if rf, ok := ret.Get(0).(func(string) []models.ActivitySuggestion); ok {
		r0 = rf(campaignDescription)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.ActivitySuggestion)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(campaignDescription)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAIService_GenerateActivitySuggestions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenerateActivitySuggestions'
type MockAIService_GenerateActivitySuggestions_Call struct {
	*mock.Call
}

// GenerateActivitySuggestions is a helper method to define mock.On call
//   - campaignDescription string
func (_e *MockAIService_Expecter) GenerateActivitySuggestions(campaignDescription interface{}) *MockAIService_GenerateActivitySuggestions_Call {
	return &MockAIService_GenerateActivitySuggestions_Call{Call: _e.mock.On("GenerateActivitySuggestions", campaignDescription)}
}

func (_c *MockAIService_GenerateActivitySuggestions_Call) Run(run func(campaignDescription string)) *MockAIService_GenerateActivitySuggestions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockAIService_GenerateActivitySuggestions_Call) Return(_a0 []models.ActivitySuggestion, _a1 error) *MockAIService_GenerateActivitySuggestions_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAIService_GenerateActivitySuggestions_Call) RunAndReturn(run func(string) ([]models.ActivitySuggestion, error)) *MockAIService_GenerateActivitySuggestions_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAIService creates a new instance of MockAIService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAIService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAIService {
	mock := &MockAIService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
