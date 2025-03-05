// Code generated by mockery v2.34.2. DO NOT EDIT.

package mocks

import (
	context "context"
	model "golang-rest-api-template/internal/handlers/category/model"

	mock "github.com/stretchr/testify/mock"

	repositorymodel "golang-rest-api-template/internal/repository/model"
)

// CategoryServiceInterface is an autogenerated mock type for the CategoryServiceInterface type
type CategoryServiceInterface struct {
	mock.Mock
}

// CreateCategory provides a mock function with given fields: ctx, _a1
func (_m *CategoryServiceInterface) CreateCategory(ctx context.Context, _a1 model.CreateCategoryRequest) (int64, error) {
	ret := _m.Called(ctx, _a1)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.CreateCategoryRequest) (int64, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.CreateCategoryRequest) int64); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.CreateCategoryRequest) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteCategory provides a mock function with given fields: ctx, id
func (_m *CategoryServiceInterface) DeleteCategory(ctx context.Context, id int) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllCategories provides a mock function with given fields: ctx
func (_m *CategoryServiceInterface) GetAllCategories(ctx context.Context) ([]repositorymodel.Category, error) {
	ret := _m.Called(ctx)

	var r0 []repositorymodel.Category
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]repositorymodel.Category, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []repositorymodel.Category); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]repositorymodel.Category)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCategoryByID provides a mock function with given fields: ctx, id
func (_m *CategoryServiceInterface) GetCategoryByID(ctx context.Context, id int) (*repositorymodel.Category, error) {
	ret := _m.Called(ctx, id)

	var r0 *repositorymodel.Category
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (*repositorymodel.Category, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) *repositorymodel.Category); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*repositorymodel.Category)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateCategory provides a mock function with given fields: ctx, id, _a2
func (_m *CategoryServiceInterface) UpdateCategory(ctx context.Context, id int, _a2 model.UpdateCategoryRequest) error {
	ret := _m.Called(ctx, id, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, model.UpdateCategoryRequest) error); ok {
		r0 = rf(ctx, id, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewCategoryServiceInterface creates a new instance of CategoryServiceInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCategoryServiceInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *CategoryServiceInterface {
	mock := &CategoryServiceInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
