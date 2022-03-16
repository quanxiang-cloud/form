package lowcode

import (
	"fmt"
	"net/http"

	"github.com/quanxiang-cloud/cabin/tailormade/client"
)

type Form interface {
	getCacheMatchRole()
	getRoleMatchPermit()
}

type form struct {
	client http.Client
}

func NewForm(conf client.Config) Form {
	return &form{
		client: client.New(conf),
	}
}

func (f *form) getCacheMatchRole() {
	fmt.Println("lowcode.getCacheMatchRole")
}

func (f *form) getRoleMatchPermit() {
	fmt.Println("lowcode.getRoleMatchPermit")
}
