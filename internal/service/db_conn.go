package service

import (
	"github.com/quanxiang-cloud/cabin/logger"
	mysql2 "github.com/quanxiang-cloud/cabin/tailormade/db/mysql"
	"github.com/quanxiang-cloud/form/pkg/misc/config"

	"gorm.io/gorm"
)

var (
	mysqlDBInst *gorm.DB
)

func CreateMysqlConn(conf *config.Config) (*gorm.DB, error) {
	if mysqlDBInst == nil {
		db, err := mysql2.New(conf.Mysql, logger.Logger)
		if err != nil {
			return nil, err
		}
		mysqlDBInst = db
	}
	return mysqlDBInst, nil
}
