package service

import "context"

type UpdateCustomResp struct {
}

type DeletePageMenuByMenuIDReq struct {
}

type DeletePageMenuByMenuIDResp struct {
}

type GetByMenuIDReq struct {
}

type GetByMenuIDResp struct {
}

type CustomPage interface {
	// CreateCustom Create CustomPage
	CreateCustom(ctx context.Context, req *CreateCustomReq) (*CreateCustomResp, error)

	// UpdateCustomPage Update customPage information
	UpdateCustomPage(ctx context.Context, req *UpdateCustomReq) (*UpdateCustomResp, error)

	// DeletePageMenuByMenuID Removes the association between the custom page and the menu
	DeletePageMenuByMenuID(ctx context.Context, req *DeletePageMenuByMenuIDReq) (*DeletePageMenuByMenuIDResp, error)

	// GetByMenuID Get the custom page information by menu id
	GetByMenuID(ctx context.Context, req *GetByMenuIDReq) (*GetByMenuIDResp, error)
}

type CreateCustomReq struct {
}

type CreateCustomResp struct {
}

type UpdateCustomReq struct {
}
