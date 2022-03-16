package lowcode

import (
	"fmt"
	"os"

	"github.com/quanxiang-cloud/cabin/tailormade/client"
)

var authURL = "%s/api/v1/jwt/check"

func init() {
	formHost := os.Getenv("FORM_HOST")
	if formHost == "" {
		formHost = "http://form"
	}
	authURL = fmt.Sprintf(authURL, formHost)
}

type Lowcode struct {
	form Form
}

func NewLowcode() *Lowcode {
	return &Lowcode{
		form: NewForm(client.Config{
			Timeout:      20,
			MaxIdleConns: 10,
		}),
	}
}

func (l *Lowcode) GetCacheMatchRole() {
	l.form.getCacheMatchRole()
}

func (l *Lowcode) GetRoleMatchPermit() {
	l.form.getRoleMatchPermit()
}
