package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/form/internal/service"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
	"net/http"
	"strings"
)

// 表单

type Form struct {
	form service.Form
}

// NewForm NewForm
func NewForm(config *config.Config) (*Form, error) {
	form, err := service.NewForm(config)
	if err != nil {
		return nil, err
	}
	return &Form{
		form: form,
	}, nil
}

// CheckURL CheckURL
func CheckURL(c *gin.Context) (appID, tableName string, err error) {
	appID, ok := c.Params.Get("appID")
	tableName, okt := c.Params.Get("tableName")
	if !ok || !okt {
		err = errors.New("invalid URI")
		return
	}
	return

}

func Search(f *Form, isAuth bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		depIDS := strings.Split(c.GetHeader(_departmentID), ",")
		req := &service.SearchReq{
			UserID: c.GetHeader(_userID),
			DepID:  depIDS[len(depIDS)-1],
			IsAuth: isAuth,
		}
		var err error
		req.AppID, req.TableID, err = CheckURL(c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		if err = c.ShouldBind(req); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		resp.Format(f.form.Search(header.MutateContext(c), req)).Context(c)
	}
}

func Get(f *Form, isAuth bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		depIDS := strings.Split(c.GetHeader(_departmentID), ",")
		req := &service.SearchReq{
			UserID: c.GetHeader(_userID),
			DepID:  depIDS[len(depIDS)-1],
			IsAuth: isAuth,
		}
		var err error
		req.AppID, req.TableID, err = CheckURL(c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		if err = c.ShouldBind(req); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		resp.Format(f.form.Search(header.MutateContext(c), req)).Context(c)

	}
}
