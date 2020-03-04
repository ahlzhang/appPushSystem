/**
 * @apidoc 巴中社保项目推送系统
 * @apiVersion -1.0.1
 * @apiBaseURL http://192.168.2.200:6002
 * @apiContent 首页
 * 开发环境地址: http://192.168.2.200:6002
 * 生产环境地址: 待定
 * 端口号 : 6002
 *
 */
package main

import (
	"google.golang.org/grpc"
	"net"
	"os"
	"sbjr.com/appPushSystem/api"
	"sbjr.com/appPushSystem/config"
	"sbjr.com/appPushSystem/control"
	pb "sbjr.com/appPushSystem/grpc"
	"sbjr.com/appPushSystem/model"
	"sbjr.com/appPushSystem/pkg/cfg"
)

func main() {
	app := config.GetApp()
	config.SetAfter(func() {
		systemRun()
	})

	app.Run(os.Args)

}

func systemRun() {
	model.StartDb()
	go control.StartPushMessageLoop()
	api.StartFunLoad()

	lis, err := net.Listen("tcp", config.Conf.System.ServicePort)
	if err != nil {
		cfg.LogFatal("推送服务无法启动,", err.Error())
	}

	s := grpc.NewServer()
	pb.RegisterPushRequestServer(s, api.PushRequestIml{})
	if err := s.Serve(lis); err != nil {
		cfg.LogFatal("启动端口监听失败:", err.Error())
	}
}
