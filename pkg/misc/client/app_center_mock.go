package client

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type AppCenterMock struct {
	mock.Mock
}

func (a AppCenterMock) CheckIsAdmin(ctx context.Context, appID, userID string, isSuper bool) (*CheckAppAdminResp, error) {
	return &CheckAppAdminResp{
		IsAdmin: true,
	}, nil
}

func NewAppCenterMock() AppCenterAPI {
	return &AppCenterMock{}
}
