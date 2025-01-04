// Code generated by mockery v2.50.0. DO NOT EDIT.

package interfaces

import (
	models "github.com/oyen-bright/goFundIt/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// MockSuggestionService is an autogenerated mock type for the SuggestionService type
type MockSuggestionService struct {
	mock.Mock
}

type MockSuggestionService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockSuggestionService) EXPECT() *MockSuggestionService_Expecter {
	return &MockSuggestionService_Expecter{mock: &_m.Mock}
}

// GetActivitySuggestions provides a mock function with given fields: campaignID, key
func (_m *MockSuggestionService) GetActivitySuggestions(campaignID string, key string) ([]models.ActivitySuggestion, error) {
	ret := _m.Called(campaignID, key)

	if len(ret) == 0 {
		panic("no return value specified for GetActivitySuggestions")
	}

	var r0 []models.ActivitySuggestion
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) ([]models.ActivitySuggestion, error)); ok {
		return rf(campaignID, key)
	}
	if rf, ok := ret.Get(0).(func(string, string) []models.ActivitySuggestion); ok {
		r0 = rf(campaignID, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.ActivitySuggestion)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(campaignID, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSuggestionService_GetActivitySuggestions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetActivitySuggestions'
type MockSuggestionService_GetActivitySuggestions_Call struct {
	*mock.Call
}

// GetActivitySuggestions is a helper method to define mock.On call
//   - campaignID string
//   - key string
func (_e *MockSuggestionService_Expecter) GetActivitySuggestions(campaignID interface{}, key interface{}) *MockSuggestionService_GetActivitySuggestions_Call {
	return &MockSuggestionService_GetActivitySuggestions_Call{Call: _e.mock.On("GetActivitySuggestions", campaignID, key)}
}

func (_c *MockSuggestionService_GetActivitySuggestions_Call) Run(run func(campaignID string, key string)) *MockSuggestionService_GetActivitySuggestions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *MockSuggestionService_GetActivitySuggestions_Call) Return(_a0 []models.ActivitySuggestion, _a1 error) *MockSuggestionService_GetActivitySuggestions_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSuggestionService_GetActivitySuggestions_Call) RunAndReturn(run func(string, string) ([]models.ActivitySuggestion, error)) *MockSuggestionService_GetActivitySuggestions_Call {
	_c.Call.Return(run)
	return _c
}

// GetActivitySuggestionsViaText provides a mock function with given fields: content
func (_m *MockSuggestionService) GetActivitySuggestionsViaText(content string) ([]models.ActivitySuggestion, error) {
	ret := _m.Called(content)

	if len(ret) == 0 {
		panic("no return value specified for GetActivitySuggestionsViaText")
	}

	var r0 []models.ActivitySuggestion
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]models.ActivitySuggestion, error)); ok {
		return rf(content)
	}
	if rf, ok := ret.Get(0).(func(string) []models.ActivitySuggestion); ok {
		r0 = rf(content)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.ActivitySuggestion)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(content)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSuggestionService_GetActivitySuggestionsViaText_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetActivitySuggestionsViaText'
type MockSuggestionService_GetActivitySuggestionsViaText_Call struct {
	*mock.Call
}

// GetActivitySuggestionsViaText is a helper method to define mock.On call
//   - content string
func (_e *MockSuggestionService_Expecter) GetActivitySuggestionsViaText(content interface{}) *MockSuggestionService_GetActivitySuggestionsViaText_Call {
	return &MockSuggestionService_GetActivitySuggestionsViaText_Call{Call: _e.mock.On("GetActivitySuggestionsViaText", content)}
}

func (_c *MockSuggestionService_GetActivitySuggestionsViaText_Call) Run(run func(content string)) *MockSuggestionService_GetActivitySuggestionsViaText_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockSuggestionService_GetActivitySuggestionsViaText_Call) Return(_a0 []models.ActivitySuggestion, _a1 error) *MockSuggestionService_GetActivitySuggestionsViaText_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSuggestionService_GetActivitySuggestionsViaText_Call) RunAndReturn(run func(string) ([]models.ActivitySuggestion, error)) *MockSuggestionService_GetActivitySuggestionsViaText_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockSuggestionService creates a new instance of MockSuggestionService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockSuggestionService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockSuggestionService {
	mock := &MockSuggestionService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
