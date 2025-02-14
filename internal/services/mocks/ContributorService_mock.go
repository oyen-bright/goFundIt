// Code generated by mockery v2.50.0. DO NOT EDIT.

package interfaces

import (
	models "github.com/oyen-bright/goFundIt/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// MockContributorService is an autogenerated mock type for the ContributorService type
type MockContributorService struct {
	mock.Mock
}

type MockContributorService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockContributorService) EXPECT() *MockContributorService_Expecter {
	return &MockContributorService_Expecter{mock: &_m.Mock}
}

// AddContributorToCampaign provides a mock function with given fields: contribution, campaignId, campaignKey, userHandle
func (_m *MockContributorService) AddContributorToCampaign(contribution *models.Contributor, campaignId string, campaignKey string, userHandle string) error {
	ret := _m.Called(contribution, campaignId, campaignKey, userHandle)

	if len(ret) == 0 {
		panic("no return value specified for AddContributorToCampaign")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Contributor, string, string, string) error); ok {
		r0 = rf(contribution, campaignId, campaignKey, userHandle)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockContributorService_AddContributorToCampaign_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddContributorToCampaign'
type MockContributorService_AddContributorToCampaign_Call struct {
	*mock.Call
}

// AddContributorToCampaign is a helper method to define mock.On call
//   - contribution *models.Contributor
//   - campaignId string
//   - campaignKey string
//   - userHandle string
func (_e *MockContributorService_Expecter) AddContributorToCampaign(contribution interface{}, campaignId interface{}, campaignKey interface{}, userHandle interface{}) *MockContributorService_AddContributorToCampaign_Call {
	return &MockContributorService_AddContributorToCampaign_Call{Call: _e.mock.On("AddContributorToCampaign", contribution, campaignId, campaignKey, userHandle)}
}

func (_c *MockContributorService_AddContributorToCampaign_Call) Run(run func(contribution *models.Contributor, campaignId string, campaignKey string, userHandle string)) *MockContributorService_AddContributorToCampaign_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*models.Contributor), args[1].(string), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *MockContributorService_AddContributorToCampaign_Call) Return(_a0 error) *MockContributorService_AddContributorToCampaign_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockContributorService_AddContributorToCampaign_Call) RunAndReturn(run func(*models.Contributor, string, string, string) error) *MockContributorService_AddContributorToCampaign_Call {
	_c.Call.Return(run)
	return _c
}

// GetContributorByID provides a mock function with given fields: contributorID
func (_m *MockContributorService) GetContributorByID(contributorID uint) (models.Contributor, error) {
	ret := _m.Called(contributorID)

	if len(ret) == 0 {
		panic("no return value specified for GetContributorByID")
	}

	var r0 models.Contributor
	var r1 error
	if rf, ok := ret.Get(0).(func(uint) (models.Contributor, error)); ok {
		return rf(contributorID)
	}
	if rf, ok := ret.Get(0).(func(uint) models.Contributor); ok {
		r0 = rf(contributorID)
	} else {
		r0 = ret.Get(0).(models.Contributor)
	}

	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(contributorID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockContributorService_GetContributorByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetContributorByID'
type MockContributorService_GetContributorByID_Call struct {
	*mock.Call
}

// GetContributorByID is a helper method to define mock.On call
//   - contributorID uint
func (_e *MockContributorService_Expecter) GetContributorByID(contributorID interface{}) *MockContributorService_GetContributorByID_Call {
	return &MockContributorService_GetContributorByID_Call{Call: _e.mock.On("GetContributorByID", contributorID)}
}

func (_c *MockContributorService_GetContributorByID_Call) Run(run func(contributorID uint)) *MockContributorService_GetContributorByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uint))
	})
	return _c
}

func (_c *MockContributorService_GetContributorByID_Call) Return(_a0 models.Contributor, _a1 error) *MockContributorService_GetContributorByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockContributorService_GetContributorByID_Call) RunAndReturn(run func(uint) (models.Contributor, error)) *MockContributorService_GetContributorByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetContributorsByCampaignID provides a mock function with given fields: campaignID
func (_m *MockContributorService) GetContributorsByCampaignID(campaignID string) ([]models.Contributor, error) {
	ret := _m.Called(campaignID)

	if len(ret) == 0 {
		panic("no return value specified for GetContributorsByCampaignID")
	}

	var r0 []models.Contributor
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]models.Contributor, error)); ok {
		return rf(campaignID)
	}
	if rf, ok := ret.Get(0).(func(string) []models.Contributor); ok {
		r0 = rf(campaignID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Contributor)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(campaignID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockContributorService_GetContributorsByCampaignID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetContributorsByCampaignID'
type MockContributorService_GetContributorsByCampaignID_Call struct {
	*mock.Call
}

// GetContributorsByCampaignID is a helper method to define mock.On call
//   - campaignID string
func (_e *MockContributorService_Expecter) GetContributorsByCampaignID(campaignID interface{}) *MockContributorService_GetContributorsByCampaignID_Call {
	return &MockContributorService_GetContributorsByCampaignID_Call{Call: _e.mock.On("GetContributorsByCampaignID", campaignID)}
}

func (_c *MockContributorService_GetContributorsByCampaignID_Call) Run(run func(campaignID string)) *MockContributorService_GetContributorsByCampaignID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockContributorService_GetContributorsByCampaignID_Call) Return(_a0 []models.Contributor, _a1 error) *MockContributorService_GetContributorsByCampaignID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockContributorService_GetContributorsByCampaignID_Call) RunAndReturn(run func(string) ([]models.Contributor, error)) *MockContributorService_GetContributorsByCampaignID_Call {
	_c.Call.Return(run)
	return _c
}

// RemoveContributorFromCampaign provides a mock function with given fields: contributorId, campaignId, userHandle, key
func (_m *MockContributorService) RemoveContributorFromCampaign(contributorId uint, campaignId string, userHandle string, key string) error {
	ret := _m.Called(contributorId, campaignId, userHandle, key)

	if len(ret) == 0 {
		panic("no return value specified for RemoveContributorFromCampaign")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uint, string, string, string) error); ok {
		r0 = rf(contributorId, campaignId, userHandle, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockContributorService_RemoveContributorFromCampaign_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoveContributorFromCampaign'
type MockContributorService_RemoveContributorFromCampaign_Call struct {
	*mock.Call
}

// RemoveContributorFromCampaign is a helper method to define mock.On call
//   - contributorId uint
//   - campaignId string
//   - userHandle string
//   - key string
func (_e *MockContributorService_Expecter) RemoveContributorFromCampaign(contributorId interface{}, campaignId interface{}, userHandle interface{}, key interface{}) *MockContributorService_RemoveContributorFromCampaign_Call {
	return &MockContributorService_RemoveContributorFromCampaign_Call{Call: _e.mock.On("RemoveContributorFromCampaign", contributorId, campaignId, userHandle, key)}
}

func (_c *MockContributorService_RemoveContributorFromCampaign_Call) Run(run func(contributorId uint, campaignId string, userHandle string, key string)) *MockContributorService_RemoveContributorFromCampaign_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uint), args[1].(string), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *MockContributorService_RemoveContributorFromCampaign_Call) Return(_a0 error) *MockContributorService_RemoveContributorFromCampaign_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockContributorService_RemoveContributorFromCampaign_Call) RunAndReturn(run func(uint, string, string, string) error) *MockContributorService_RemoveContributorFromCampaign_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateContributor provides a mock function with given fields: contributor
func (_m *MockContributorService) UpdateContributor(contributor *models.Contributor) error {
	ret := _m.Called(contributor)

	if len(ret) == 0 {
		panic("no return value specified for UpdateContributor")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Contributor) error); ok {
		r0 = rf(contributor)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockContributorService_UpdateContributor_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateContributor'
type MockContributorService_UpdateContributor_Call struct {
	*mock.Call
}

// UpdateContributor is a helper method to define mock.On call
//   - contributor *models.Contributor
func (_e *MockContributorService_Expecter) UpdateContributor(contributor interface{}) *MockContributorService_UpdateContributor_Call {
	return &MockContributorService_UpdateContributor_Call{Call: _e.mock.On("UpdateContributor", contributor)}
}

func (_c *MockContributorService_UpdateContributor_Call) Run(run func(contributor *models.Contributor)) *MockContributorService_UpdateContributor_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*models.Contributor))
	})
	return _c
}

func (_c *MockContributorService_UpdateContributor_Call) Return(_a0 error) *MockContributorService_UpdateContributor_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockContributorService_UpdateContributor_Call) RunAndReturn(run func(*models.Contributor) error) *MockContributorService_UpdateContributor_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateContributorByID provides a mock function with given fields: contributor, contributorID, userEmail
func (_m *MockContributorService) UpdateContributorByID(contributor *models.Contributor, contributorID uint, userEmail string) (models.Contributor, error) {
	ret := _m.Called(contributor, contributorID, userEmail)

	if len(ret) == 0 {
		panic("no return value specified for UpdateContributorByID")
	}

	var r0 models.Contributor
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.Contributor, uint, string) (models.Contributor, error)); ok {
		return rf(contributor, contributorID, userEmail)
	}
	if rf, ok := ret.Get(0).(func(*models.Contributor, uint, string) models.Contributor); ok {
		r0 = rf(contributor, contributorID, userEmail)
	} else {
		r0 = ret.Get(0).(models.Contributor)
	}

	if rf, ok := ret.Get(1).(func(*models.Contributor, uint, string) error); ok {
		r1 = rf(contributor, contributorID, userEmail)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockContributorService_UpdateContributorByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateContributorByID'
type MockContributorService_UpdateContributorByID_Call struct {
	*mock.Call
}

// UpdateContributorByID is a helper method to define mock.On call
//   - contributor *models.Contributor
//   - contributorID uint
//   - userEmail string
func (_e *MockContributorService_Expecter) UpdateContributorByID(contributor interface{}, contributorID interface{}, userEmail interface{}) *MockContributorService_UpdateContributorByID_Call {
	return &MockContributorService_UpdateContributorByID_Call{Call: _e.mock.On("UpdateContributorByID", contributor, contributorID, userEmail)}
}

func (_c *MockContributorService_UpdateContributorByID_Call) Run(run func(contributor *models.Contributor, contributorID uint, userEmail string)) *MockContributorService_UpdateContributorByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*models.Contributor), args[1].(uint), args[2].(string))
	})
	return _c
}

func (_c *MockContributorService_UpdateContributorByID_Call) Return(retrievedContributor models.Contributor, err error) *MockContributorService_UpdateContributorByID_Call {
	_c.Call.Return(retrievedContributor, err)
	return _c
}

func (_c *MockContributorService_UpdateContributorByID_Call) RunAndReturn(run func(*models.Contributor, uint, string) (models.Contributor, error)) *MockContributorService_UpdateContributorByID_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockContributorService creates a new instance of MockContributorService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockContributorService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockContributorService {
	mock := &MockContributorService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
