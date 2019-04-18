package model

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"jiaotou.com/appPushSystem/config"
	"jiaotou.com/appPushSystem/pkg/cfg"
)

var MysqlDefault *xorm.Engine

func StartDb() {
	cfg.LogInfo("数据库初始化......")

	var err error

	MysqlDefault, err = xorm.NewEngine("mysql", config.Conf.Sql.WriteConn)
	if err != nil {
		cfg.LogErr("数据库连接失败:", err.Error())
		return
	}

	MysqlDefault.ShowSQL(config.IsTest)

	MysqlDefault.SetMaxOpenConns(config.Conf.Sql.MaxConn)
	MysqlDefault.SetMaxIdleConns(config.Conf.Sql.MaxIdleConn)
}

func GetDb() *xorm.Engine {
	return MysqlDefault
}
