// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	avito_testcase "avito_testcase"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAuthorization is a mock of Authorization interface.
type MockAuthorization struct {
	ctrl     *gomock.Controller
	recorder *MockAuthorizationMockRecorder
}

// MockAuthorizationMockRecorder is the mock recorder for MockAuthorization.
type MockAuthorizationMockRecorder struct {
	mock *MockAuthorization
}

// NewMockAuthorization creates a new mock instance.
func NewMockAuthorization(ctrl *gomock.Controller) *MockAuthorization {
	mock := &MockAuthorization{ctrl: ctrl}
	mock.recorder = &MockAuthorizationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthorization) EXPECT() *MockAuthorizationMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockAuthorization) CreateUser(user avito_testcase.User) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", user)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockAuthorizationMockRecorder) CreateUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockAuthorization)(nil).CreateUser), user)
}

// GenerateToken mocks base method.
func (m *MockAuthorization) GenerateToken(username, password string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateToken", username, password)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MockAuthorizationMockRecorder) GenerateToken(username, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockAuthorization)(nil).GenerateToken), username, password)
}

// GetUserStatus mocks base method.
func (m *MockAuthorization) GetUserStatus(id int) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserStatus", id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserStatus indicates an expected call of GetUserStatus.
func (mr *MockAuthorizationMockRecorder) GetUserStatus(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserStatus", reflect.TypeOf((*MockAuthorization)(nil).GetUserStatus), id)
}

// ParseToken mocks base method.
func (m *MockAuthorization) ParseToken(token string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseToken", token)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseToken indicates an expected call of ParseToken.
func (mr *MockAuthorizationMockRecorder) ParseToken(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseToken", reflect.TypeOf((*MockAuthorization)(nil).ParseToken), token)
}

// MockBanner is a mock of Banner interface.
type MockBanner struct {
	ctrl     *gomock.Controller
	recorder *MockBannerMockRecorder
}

// MockBannerMockRecorder is the mock recorder for MockBanner.
type MockBannerMockRecorder struct {
	mock *MockBanner
}

// NewMockBanner creates a new mock instance.
func NewMockBanner(ctrl *gomock.Controller) *MockBanner {
	mock := &MockBanner{ctrl: ctrl}
	mock.recorder = &MockBannerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBanner) EXPECT() *MockBannerMockRecorder {
	return m.recorder
}

// CreateBanner mocks base method.
func (m *MockBanner) CreateBanner(banner avito_testcase.Banners, tagIDs []int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBanner", banner, tagIDs)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBanner indicates an expected call of CreateBanner.
func (mr *MockBannerMockRecorder) CreateBanner(banner, tagIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBanner", reflect.TypeOf((*MockBanner)(nil).CreateBanner), banner, tagIDs)
}

// DeleteBanner mocks base method.
func (m *MockBanner) DeleteBanner(bannerID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBanner", bannerID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBanner indicates an expected call of DeleteBanner.
func (mr *MockBannerMockRecorder) DeleteBanner(bannerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBanner", reflect.TypeOf((*MockBanner)(nil).DeleteBanner), bannerID)
}

// GetAllBanners mocks base method.
func (m *MockBanner) GetAllBanners(featureID, tagID, limit, offset int) ([]avito_testcase.AllBanners, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllBanners", featureID, tagID, limit, offset)
	ret0, _ := ret[0].([]avito_testcase.AllBanners)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllBanners indicates an expected call of GetAllBanners.
func (mr *MockBannerMockRecorder) GetAllBanners(featureID, tagID, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllBanners", reflect.TypeOf((*MockBanner)(nil).GetAllBanners), featureID, tagID, limit, offset)
}

// GetAllBannersForAdmin mocks base method.
func (m *MockBanner) GetAllBannersForAdmin(featureID, tagID, limit, offset int) ([]avito_testcase.AllBanners, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllBannersForAdmin", featureID, tagID, limit, offset)
	ret0, _ := ret[0].([]avito_testcase.AllBanners)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllBannersForAdmin indicates an expected call of GetAllBannersForAdmin.
func (mr *MockBannerMockRecorder) GetAllBannersForAdmin(featureID, tagID, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllBannersForAdmin", reflect.TypeOf((*MockBanner)(nil).GetAllBannersForAdmin), featureID, tagID, limit, offset)
}

// GetBannerByTagAndFeature mocks base method.
func (m *MockBanner) GetBannerByTagAndFeature(tagID, featureID int, lastVersion bool) (avito_testcase.Banners, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBannerByTagAndFeature", tagID, featureID, lastVersion)
	ret0, _ := ret[0].(avito_testcase.Banners)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBannerByTagAndFeature indicates an expected call of GetBannerByTagAndFeature.
func (mr *MockBannerMockRecorder) GetBannerByTagAndFeature(tagID, featureID, lastVersion interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBannerByTagAndFeature", reflect.TypeOf((*MockBanner)(nil).GetBannerByTagAndFeature), tagID, featureID, lastVersion)
}

// GetBannerByTagAndFeatureForAdmin mocks base method.
func (m *MockBanner) GetBannerByTagAndFeatureForAdmin(tagID, featureID int, lastVersion bool) (avito_testcase.Banners, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBannerByTagAndFeatureForAdmin", tagID, featureID, lastVersion)
	ret0, _ := ret[0].(avito_testcase.Banners)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBannerByTagAndFeatureForAdmin indicates an expected call of GetBannerByTagAndFeatureForAdmin.
func (mr *MockBannerMockRecorder) GetBannerByTagAndFeatureForAdmin(tagID, featureID, lastVersion interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBannerByTagAndFeatureForAdmin", reflect.TypeOf((*MockBanner)(nil).GetBannerByTagAndFeatureForAdmin), tagID, featureID, lastVersion)
}

// UpdateBanner mocks base method.
func (m *MockBanner) UpdateBanner(bannerID int, input avito_testcase.UpdateBanner) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBanner", bannerID, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBanner indicates an expected call of UpdateBanner.
func (mr *MockBannerMockRecorder) UpdateBanner(bannerID, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBanner", reflect.TypeOf((*MockBanner)(nil).UpdateBanner), bannerID, input)
}