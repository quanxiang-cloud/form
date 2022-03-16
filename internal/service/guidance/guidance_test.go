package guidance

import (
	"context"
	"fmt"
	"testing"

	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/service"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CertifierSuite struct {
	suite.Suite
	ctx      context.Context
	AppID    string
	TableID  string
	permit   service.Permit
	conf     *config.Config
	UserID   string
	UserName string
}

func TestCertifier(t *testing.T) {
	suite.Run(t, new(CertifierSuite))
}

func (suite *CertifierSuite) SetupTest() {
	fmt.Println("a")
	suite.AppID = "app01"
	suite.TableID = "table01"
	suite.UserID = "userID"
	suite.UserName = "周慧婷"
	suite.ctx = context.TODO()
	var err error
	suite.conf, err = config.NewConfig("../../../configs/config.yml")
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), suite.conf)

	suite.permit, err = service.NewPermit(suite.conf)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), suite.conf)
}

func (suite *CertifierSuite) CertifierBefore() {
}

func (suite *CertifierSuite) CertifierAfter() {
}

func (suite *CertifierSuite) TestCertifier() {
	createRole := &service.CreateRoleReq{
		UserID:      suite.UserID,
		UserName:    suite.UserName,
		AppID:       suite.AppID,
		Name:        "测试角色",
		Description: "测试角色描述",
		Types:       models.CreateType,
	}
	role, err := suite.permit.CreateRole(suite.ctx, createRole)

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), role)
	// 往角色加人

	//	Authorizes []*Owners `json:"authorizes"`
	//	RoleID     string    `json:"roleID"`
	owners := &service.Owners{
		Owner:     suite.UserID,
		OwnerName: suite.UserName,
		Types:     1,
	}

	addReq := &service.AddOwnerToRoleReq{
		RoleID:     role.ID,
		Authorizes: []*service.Owners{owners},
	}

	toRole, err := suite.permit.AddOwnerToRole(suite.ctx, addReq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), toRole)
	req := &service.CreatePerReq{
		Params: models.FiledPermit{
			"field_uHU7doso": models.Key{
				Type: "string",
			},
			"field_lO0xgd9D": models.Key{
				Type: "string",
			},
		},
		Response: models.FiledPermit{
			"field_uHU7doso": models.Key{
				Type: "string",
			},
			"field_lO0xgd9D": models.Key{
				Type: "string",
			},
		},
		RoleID:   role.ID,
		UserID:   suite.UserID,
		UserName: suite.UserName,
	}
	permit, err := suite.permit.CreatePermit(suite.ctx, req)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), permit)
}
