package client

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type AppCenterMock struct {
	mock.Mock
}

func (a AppCenterMock) GetOne(ctx context.Context, appID string) (*AppResp, error) {
	return &AppResp{}, nil
}

func (a AppCenterMock) CheckIsAdmin(ctx context.Context, appID, userID string, isSuper bool) (*CheckAppAdminResp, error) {
	return &CheckAppAdminResp{
		IsAdmin: true,
	}, nil
}

func NewAppCenterMock() AppCenterAPI {
	return &AppCenterMock{}
}
