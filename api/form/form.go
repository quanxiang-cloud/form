package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/form/internal/service/consensus"
)

type profile struct {
	userID   string
	depID    string
	userName string
}

func action(ctr consensus.Guidance) gin.HandlerFunc {
	return func(c *gin.Context) {
		bus := &consensus.Bus{}
		profiles := getProfile(c)
		bus.UserID = profiles.userID
		bus.UserName = profiles.userName
		bus.DepID = profiles.depID
		bus.Method = c.Param("action")
		bus.Path = c.Request.RequestURI
		var err error
		bus.AppID, bus.TableID, err = checkURL(c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if err = c.ShouldBind(bus); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		do, err := ctr.Do(header.MutateContext(c), bus)

		resp.Format(do, err).Context(c)
	}
}

func getProfile(c *gin.Context) *profile {
	userID := c.GetHeader(_userID)
	userName := c.GetHeader(_userName)
	depIDS := strings.Split(c.GetHeader(_departmentID), ",")
	return &profile{
		userID:   userID,
		userName: userName,
		depID:    depIDS[len(depIDS)-1],
	}
}

// checkURL CheckURL
func checkURL(c *gin.Context) (appID, tableName string, err error) {
	appID, ok := c.Params.Get("appID")
	tableName, okt := c.Params.Get("tableName")
	if !ok || !okt {
		err = errors.New("invalid URI")
		return
	}
	return
}
