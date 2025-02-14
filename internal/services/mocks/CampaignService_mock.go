// Code generated by mockery v2.50.0. DO NOT EDIT.

package interfaces

import (
	dto "github.com/oyen-bright/goFundIt/internal/api/dto/campaign"

	mock "github.com/stretchr/testify/mock"

	models "github.com/oyen-bright/goFundIt/internal/models"
)

// MockCampaignService is an autogenerated mock type for the CampaignService type
type MockCampaignService struct {
	mock.Mock
}

type MockCampaignService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCampaignService) EXPECT() *MockCampaignService_Expecter {
	return &MockCampaignService_Expecter{mock: &_m.Mock}
}

// CreateCampaign provides a mock function with given fields: campaign, userHandle
func (_m *MockCampaignService) CreateCampaign(campaign *models.Campaign, userHandle string) (models.Campaign, error) {
	ret := _m.Called(campaign, userHandle)

	if len(ret) == 0 {
		panic("no return value specified for CreateCampaign")
	}

	var r0 models.Campaign
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.Campaign, string) (models.Campaign, error)); ok {
		return rf(campaign, userHandle)
	}
	if rf, ok := ret.Get(0).(func(*models.Campaign, string) models.Campaign); ok {
		r0 = rf(campaign, userHandle)
	} else {
		r0 = ret.Get(0).(models.Campaign)
	}

	if rf, ok := ret.Get(1).(func(*models.Campaign, string) error); ok {
		r1 = rf(campaign, userHandle)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCampaignService_CreateCampaign_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateCampaign'
type MockCampaignService_CreateCampaign_Call struct {
	*mock.Call
}

// CreateCampaign is a helper method to define mock.On call
//   - campaign *models.Campaign
//   - userHandle string
func (_e *MockCampaignService_Expecter) CreateCampaign(campaign interface{}, userHandle interface{}) *MockCampaignService_CreateCampaign_Call {
	return &MockCampaignService_CreateCampaign_Call{Call: _e.mock.On("CreateCampaign", campaign, userHandle)}
}

func (_c *MockCampaignService_CreateCampaign_Call) Run(run func(campaign *models.Campaign, userHandle string)) *MockCampaignService_CreateCampaign_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*models.Campaign), args[1].(string))
	})
	return _c
}

func (_c *MockCampaignService_CreateCampaign_Call) Return(_a0 models.Campaign, _a1 error) *MockCampaignService_CreateCampaign_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCampaignService_CreateCampaign_Call) RunAndReturn(run func(*models.Campaign, string) (models.Campaign, error)) *MockCampaignService_CreateCampaign_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteCampaign provides a mock function with given fields: campaignID
func (_m *MockCampaignService) DeleteCampaign(campaignID string) error {
	ret := _m.Called(campaignID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteCampaign")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(campaignID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockCampaignService_DeleteCampaign_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteCampaign'
type MockCampaignService_DeleteCampaign_Call struct {
	*mock.Call
}

// DeleteCampaign is a helper method to define mock.On call
//   - campaignID string
func (_e *MockCampaignService_Expecter) DeleteCampaign(campaignID interface{}) *MockCampaignService_DeleteCampaign_Call {
	return &MockCampaignService_DeleteCampaign_Call{Call: _e.mock.On("DeleteCampaign", campaignID)}
}

func (_c *MockCampaignService_DeleteCampaign_Call) Run(run func(campaignID string)) *MockCampaignService_DeleteCampaign_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockCampaignService_DeleteCampaign_Call) Return(_a0 error) *MockCampaignService_DeleteCampaign_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCampaignService_DeleteCampaign_Call) RunAndReturn(run func(string) error) *MockCampaignService_DeleteCampaign_Call {
	_c.Call.Return(run)
	return _c
}

// GetActiveCampaigns provides a mock function with no fields
func (_m *MockCampaignService) GetActiveCampaigns() ([]models.Campaign, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetActiveCampaigns")
	}

	var r0 []models.Campaign
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]models.Campaign, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []models.Campaign); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Campaign)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCampaignService_GetActiveCampaigns_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetActiveCampaigns'
type MockCampaignService_GetActiveCampaigns_Call struct {
	*mock.Call
}

// GetActiveCampaigns is a helper method to define mock.On call
func (_e *MockCampaignService_Expecter) GetActiveCampaigns() *MockCampaignService_GetActiveCampaigns_Call {
	return &MockCampaignService_GetActiveCampaigns_Call{Call: _e.mock.On("GetActiveCampaigns")}
}

func (_c *MockCampaignService_GetActiveCampaigns_Call) Run(run func()) *MockCampaignService_GetActiveCampaigns_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCampaignService_GetActiveCampaigns_Call) Return(_a0 []models.Campaign, _a1 error) *MockCampaignService_GetActiveCampaigns_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCampaignService_GetActiveCampaigns_Call) RunAndReturn(run func() ([]models.Campaign, error)) *MockCampaignService_GetActiveCampaigns_Call {
	_c.Call.Return(run)
	return _c
}

// GetCampaignByID provides a mock function with given fields: id, key
func (_m *MockCampaignService) GetCampaignByID(id string, key string) (*models.Campaign, error) {
	ret := _m.Called(id, key)

	if len(ret) == 0 {
		panic("no return value specified for GetCampaignByID")
	}

	var r0 *models.Campaign
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (*models.Campaign, error)); ok {
		return rf(id, key)
	}
	if rf, ok := ret.Get(0).(func(string, string) *models.Campaign); ok {
		r0 = rf(id, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Campaign)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(id, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCampaignService_GetCampaignByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCampaignByID'
type MockCampaignService_GetCampaignByID_Call struct {
	*mock.Call
}

// GetCampaignByID is a helper method to define mock.On call
//   - id string
//   - key string
func (_e *MockCampaignService_Expecter) GetCampaignByID(id interface{}, key interface{}) *MockCampaignService_GetCampaignByID_Call {
	return &MockCampaignService_GetCampaignByID_Call{Call: _e.mock.On("GetCampaignByID", id, key)}
}

func (_c *MockCampaignService_GetCampaignByID_Call) Run(run func(id string, key string)) *MockCampaignService_GetCampaignByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *MockCampaignService_GetCampaignByID_Call) Return(_a0 *models.Campaign, _a1 error) *MockCampaignService_GetCampaignByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCampaignService_GetCampaignByID_Call) RunAndReturn(run func(string, string) (*models.Campaign, error)) *MockCampaignService_GetCampaignByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetCampaignByIDWithAllRelatedData provides a mock function with given fields: id
func (_m *MockCampaignService) GetCampaignByIDWithAllRelatedData(id string) (*models.Campaign, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetCampaignByIDWithAllRelatedData")
	}

	var r0 *models.Campaign
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*models.Campaign, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *models.Campaign); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Campaign)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCampaignService_GetCampaignByIDWithAllRelatedData_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCampaignByIDWithAllRelatedData'
type MockCampaignService_GetCampaignByIDWithAllRelatedData_Call struct {
	*mock.Call
}

// GetCampaignByIDWithAllRelatedData is a helper method to define mock.On call
//   - id string
func (_e *MockCampaignService_Expecter) GetCampaignByIDWithAllRelatedData(id interface{}) *MockCampaignService_GetCampaignByIDWithAllRelatedData_Call {
	return &MockCampaignService_GetCampaignByIDWithAllRelatedData_Call{Call: _e.mock.On("GetCampaignByIDWithAllRelatedData", id)}
}

func (_c *MockCampaignService_GetCampaignByIDWithAllRelatedData_Call) Run(run func(id string)) *MockCampaignService_GetCampaignByIDWithAllRelatedData_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockCampaignService_GetCampaignByIDWithAllRelatedData_Call) Return(_a0 *models.Campaign, _a1 error) *MockCampaignService_GetCampaignByIDWithAllRelatedData_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCampaignService_GetCampaignByIDWithAllRelatedData_Call) RunAndReturn(run func(string) (*models.Campaign, error)) *MockCampaignService_GetCampaignByIDWithAllRelatedData_Call {
	_c.Call.Return(run)
	return _c
}

// GetCampaignByIDWithContributors provides a mock function with given fields: id
func (_m *MockCampaignService) GetCampaignByIDWithContributors(id string) (*models.Campaign, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetCampaignByIDWithContributors")
	}

	var r0 *models.Campaign
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*models.Campaign, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *models.Campaign); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Campaign)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCampaignService_GetCampaignByIDWithContributors_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCampaignByIDWithContributors'
type MockCampaignService_GetCampaignByIDWithContributors_Call struct {
	*mock.Call
}

// GetCampaignByIDWithContributors is a helper method to define mock.On call
//   - id string
func (_e *MockCampaignService_Expecter) GetCampaignByIDWithContributors(id interface{}) *MockCampaignService_GetCampaignByIDWithContributors_Call {
	return &MockCampaignService_GetCampaignByIDWithContributors_Call{Call: _e.mock.On("GetCampaignByIDWithContributors", id)}
}

func (_c *MockCampaignService_GetCampaignByIDWithContributors_Call) Run(run func(id string)) *MockCampaignService_GetCampaignByIDWithContributors_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockCampaignService_GetCampaignByIDWithContributors_Call) Return(_a0 *models.Campaign, _a1 error) *MockCampaignService_GetCampaignByIDWithContributors_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCampaignService_GetCampaignByIDWithContributors_Call) RunAndReturn(run func(string) (*models.Campaign, error)) *MockCampaignService_GetCampaignByIDWithContributors_Call {
	_c.Call.Return(run)
	return _c
}

// GetExpiredCampaigns provides a mock function with no fields
func (_m *MockCampaignService) GetExpiredCampaigns() ([]models.Campaign, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetExpiredCampaigns")
	}

	var r0 []models.Campaign
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]models.Campaign, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []models.Campaign); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Campaign)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCampaignService_GetExpiredCampaigns_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetExpiredCampaigns'
type MockCampaignService_GetExpiredCampaigns_Call struct {
	*mock.Call
}

// GetExpiredCampaigns is a helper method to define mock.On call
func (_e *MockCampaignService_Expecter) GetExpiredCampaigns() *MockCampaignService_GetExpiredCampaigns_Call {
	return &MockCampaignService_GetExpiredCampaigns_Call{Call: _e.mock.On("GetExpiredCampaigns")}
}

func (_c *MockCampaignService_GetExpiredCampaigns_Call) Run(run func()) *MockCampaignService_GetExpiredCampaigns_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCampaignService_GetExpiredCampaigns_Call) Return(_a0 []models.Campaign, _a1 error) *MockCampaignService_GetExpiredCampaigns_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCampaignService_GetExpiredCampaigns_Call) RunAndReturn(run func() ([]models.Campaign, error)) *MockCampaignService_GetExpiredCampaigns_Call {
	_c.Call.Return(run)
	return _c
}

// GetNearEndCampaigns provides a mock function with no fields
func (_m *MockCampaignService) GetNearEndCampaigns() ([]models.Campaign, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetNearEndCampaigns")
	}

	var r0 []models.Campaign
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]models.Campaign, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []models.Campaign); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Campaign)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCampaignService_GetNearEndCampaigns_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetNearEndCampaigns'
type MockCampaignService_GetNearEndCampaigns_Call struct {
	*mock.Call
}

// GetNearEndCampaigns is a helper method to define mock.On call
func (_e *MockCampaignService_Expecter) GetNearEndCampaigns() *MockCampaignService_GetNearEndCampaigns_Call {
	return &MockCampaignService_GetNearEndCampaigns_Call{Call: _e.mock.On("GetNearEndCampaigns")}
}

func (_c *MockCampaignService_GetNearEndCampaigns_Call) Run(run func()) *MockCampaignService_GetNearEndCampaigns_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCampaignService_GetNearEndCampaigns_Call) Return(_a0 []models.Campaign, _a1 error) *MockCampaignService_GetNearEndCampaigns_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCampaignService_GetNearEndCampaigns_Call) RunAndReturn(run func() ([]models.Campaign, error)) *MockCampaignService_GetNearEndCampaigns_Call {
	_c.Call.Return(run)
	return _c
}

// RecalculateTargetAmount provides a mock function with given fields: campaignID
func (_m *MockCampaignService) RecalculateTargetAmount(campaignID string) {
	_m.Called(campaignID)
}

// MockCampaignService_RecalculateTargetAmount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RecalculateTargetAmount'
type MockCampaignService_RecalculateTargetAmount_Call struct {
	*mock.Call
}

// RecalculateTargetAmount is a helper method to define mock.On call
//   - campaignID string
func (_e *MockCampaignService_Expecter) RecalculateTargetAmount(campaignID interface{}) *MockCampaignService_RecalculateTargetAmount_Call {
	return &MockCampaignService_RecalculateTargetAmount_Call{Call: _e.mock.On("RecalculateTargetAmount", campaignID)}
}

func (_c *MockCampaignService_RecalculateTargetAmount_Call) Run(run func(campaignID string)) *MockCampaignService_RecalculateTargetAmount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockCampaignService_RecalculateTargetAmount_Call) Return() *MockCampaignService_RecalculateTargetAmount_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockCampaignService_RecalculateTargetAmount_Call) RunAndReturn(run func(string)) *MockCampaignService_RecalculateTargetAmount_Call {
	_c.Run(run)
	return _c
}

// UpdateCampaign provides a mock function with given fields: data, campaignID, key, userHandle
func (_m *MockCampaignService) UpdateCampaign(data dto.CampaignUpdateRequest, campaignID string, key string, userHandle string) (*models.Campaign, error) {
	ret := _m.Called(data, campaignID, key, userHandle)

	if len(ret) == 0 {
		panic("no return value specified for UpdateCampaign")
	}

	var r0 *models.Campaign
	var r1 error
	if rf, ok := ret.Get(0).(func(dto.CampaignUpdateRequest, string, string, string) (*models.Campaign, error)); ok {
		return rf(data, campaignID, key, userHandle)
	}
	if rf, ok := ret.Get(0).(func(dto.CampaignUpdateRequest, string, string, string) *models.Campaign); ok {
		r0 = rf(data, campaignID, key, userHandle)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Campaign)
		}
	}

	if rf, ok := ret.Get(1).(func(dto.CampaignUpdateRequest, string, string, string) error); ok {
		r1 = rf(data, campaignID, key, userHandle)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCampaignService_UpdateCampaign_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateCampaign'
type MockCampaignService_UpdateCampaign_Call struct {
	*mock.Call
}

// UpdateCampaign is a helper method to define mock.On call
//   - data dto.CampaignUpdateRequest
//   - campaignID string
//   - key string
//   - userHandle string
func (_e *MockCampaignService_Expecter) UpdateCampaign(data interface{}, campaignID interface{}, key interface{}, userHandle interface{}) *MockCampaignService_UpdateCampaign_Call {
	return &MockCampaignService_UpdateCampaign_Call{Call: _e.mock.On("UpdateCampaign", data, campaignID, key, userHandle)}
}

func (_c *MockCampaignService_UpdateCampaign_Call) Run(run func(data dto.CampaignUpdateRequest, campaignID string, key string, userHandle string)) *MockCampaignService_UpdateCampaign_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(dto.CampaignUpdateRequest), args[1].(string), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *MockCampaignService_UpdateCampaign_Call) Return(_a0 *models.Campaign, _a1 error) *MockCampaignService_UpdateCampaign_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCampaignService_UpdateCampaign_Call) RunAndReturn(run func(dto.CampaignUpdateRequest, string, string, string) (*models.Campaign, error)) *MockCampaignService_UpdateCampaign_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockCampaignService creates a new instance of MockCampaignService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCampaignService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCampaignService {
	mock := &MockCampaignService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
