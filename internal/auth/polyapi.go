package auth

import (
	"context"
)

type PolyAuth interface {
	Auth(context.Context, *PolyAuthReq) (*PolyAuthResp, error)
}
type polyAuth struct{}

func NewPolyAuth() PolyAuth {
	return &polyAuth{}
}

type PolyAuthReq struct{}

type PolyAuthResp struct {
	IsPermit bool
}

func (p *polyAuth) Auth(context.Context, *PolyAuthReq) (*PolyAuthResp, error) {
	// TODO: implement poly auth
	return &PolyAuthResp{IsPermit: true}, nil
}
