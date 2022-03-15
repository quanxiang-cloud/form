package lowcode

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PolyAuth struct{}

func NewPolyAuth() *PolyAuth {
	return &PolyAuth{}
}

func (p *PolyAuth) Auth(ctx *gin.Context) bool {
	// TODO: implement poly auth
	return true
}

func (p *PolyAuth) Filter(resp *http.Response) error {
	return nil
}
