/**
 * 配置文件读取.
 *
 * User: zhangbob
 * Date: 2018-12-22
 * Time: 15:15
 */
package config

import (
	"github.com/BurntSushi/toml"
	"gopkg.in/urfave/cli.v2"
	"jiaotou.com/appPushSystem/pkg/cfg"
	"jiaotou.com/appPushSystem/utils"
)

var Conf Config
var App *cli.App

var IsTest bool
var AppConfPath string

func GetApp() *cli.App {
	App = &cli.App{
		Name:    "jtSystem",
		Version: "1.0.0",
		Usage:   "交投项目-api系统",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "dir",
				Aliases: []string{"d"},
				Value:   utils.GetCurrentDirectory(),
				Usage:   "运行目录",
			},
			&cli.BoolFlag{
				Name:    "is_test",
				Aliases: []string{"i"},
				Value:   false,
				Usage:   "是否测试环境",
			},
		},
		Action: action,
	}

	return App
}

func action(c *cli.Context) error {
	AppConfPath = c.String("dir")
	IsTest = c.Bool("is_test")

	cfg.InitLogger(AppConfPath+"/conf/log/", "jtSystem.log", 51200)
	readConfig(AppConfPath + "/conf/app.toml")

	cfg.InitLogLevel(Conf.Log.LogLevel)

	return nil
}

func SetAfter(f func()) {
	App.After = func(context *cli.Context) error {
		f()
		return nil
	}
}

func readConfig(dir string) Config {
	if _, err := toml.DecodeFile(dir, &Conf); err != nil {
		cfg.LogFatal("配置文件有误:", err.Error())
	}

	cfg.LogInfof("配置文件:%+v", Conf)

	return Conf
}

type Config struct {
	System   system
	Log      log
	Sql      sql
	PushInfo pushInfo
}

type system struct {
	ServicePort string
}

type log struct {
	LogLevel int
	LogFlag  int
}

type sql struct {
	WriteConn   string
	MaxConn     int
	MaxIdleConn int
}

type pushInfo struct {
	Android android
	Ios     iOS
}

type android struct {
	AppId        string
	AppKey       string
	MasterSecret string
}

type iOS struct {
	Model    string
	Password string
	Pem      string
	Package  string
}
