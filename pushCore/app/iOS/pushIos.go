/**
 * iOS推送.
 *
 * User: zhangbob
 * Date: 2018/5/9
 * Time: 上午10:16
 */
package iOS

import (
	"errors"
	"jiaotou.com/appPushSystem/config"
	"jiaotou.com/appPushSystem/pkg/cfg"

	apns "github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"jiaotou.com/appPushSystem/pushCore"
	"sync"
)

var once sync.Once
var instance messageHandle

func GetIosHandleInstance() messageHandle {
	once.Do(func() {
		cert, err := certificate.FromPemFile(config.AppConfPath+"/conf/"+config.Conf.PushInfo.Ios.Pem, config.Conf.PushInfo.Ios.Password)
		if err != nil {
			cfg.LogFatal("ios证书验证不正确:", err)
		}

		instance.packageName = config.Conf.PushInfo.Ios.Package

		if instance.packageName == "" {
			cfg.LogFatal("ios包名未配置。")
		}

		if !config.Conf.PushInfo.Ios.IsDev {
			instance.apnsClient = apns.NewClient(cert).Production()
		} else {
			instance.apnsClient = apns.NewClient(cert).Development()
		}
	})

	return instance
}

type messageHandle struct {
	apnsClient  *apns.Client
	packageName string
}

func (t messageHandle) PushSingle(message pushCore.IMessage, callBack pushCore.IHandleMessageCallback) {
	_, clientId, content := message.ToMessage()
	notification := &apns.Notification{}
	notification.DeviceToken = clientId
	notification.Payload = []byte(content)
	notification.Topic = t.packageName

	res, err := t.apnsClient.Push(notification)
	if err != nil {
		cfg.LogErrf("iOS推送失败:%s,返回值:%+v", err.Error(), res)
		callBack.Fail(message, err)
		return
	}

	// 200  Success
	// 400 Bad request
	// 403 There was an error with the certificate.
	// 405 The request used a bad :method value. Only POST requests are supported.
	// 410 The device token is no longer active for the topic.
	// 413 The notification payload was too large.
	// 429 The server received too many requests for the same device token.
	// 500 Internal server error
	// 503 The server is shutting down and unavailable.
	if res.Sent() {
		callBack.Success(message)
	} else {
		cfg.LogErrf("iOS推失败,code:%d,原因:%s,推送token:%s", res.StatusCode, res.Reason, notification.DeviceToken)
		callBack.Fail(message, errors.New(res.Reason))
	}
}
