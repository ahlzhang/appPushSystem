/**
 * iOS推送.
 *
 * User: zhangbob
 * Date: 2018/5/9
 * Time: 上午10:16
 */
package control

import (
	"encoding/json"
	apns "github.com/sideshow/apns2"
	"jiaotou.com/appPushSystem/config"
	"jiaotou.com/appPushSystem/model"
	"jiaotou.com/appPushSystem/pkg/cfg"
)

var apnsClient *apns.Client

var iosPushList = make(chan ListSendMsg)

func sendLoopIos() {
	for {
		if v, ok := <-iosPushList; ok {
			PushIos(v.PushId, v.Token, v.Content)
		}
	}
}

func push(pushId int64, notification *apns.Notification) error {
	res, err := apnsClient.Push(notification)
	if err != nil {
		cfg.LogErr("error:", err, " res:", res)

		return err
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
		cfg.LogInfo("APNs ID:", res.ApnsID)

		model.UpdatePushBySingle(pushId, config.PushStateAlready)
	} else {
		cfg.LogErr("push error,  code :", res.StatusCode, " reason:", res.Reason, " token:", notification.DeviceToken)
		model.UpdatePushBySingle(pushId, config.PushStateFail)
	}
	return nil
}

func PushIos(pushId int64, token string, content string) {
	notification := &apns.Notification{}
	notification.DeviceToken = token
	notification.Payload = []byte(content)
	notification.Topic = config.Conf.PushInfo.Ios.Package
	push(pushId, notification)
}

func iosEncode(id int64, types int64, description string, title string) string {
	out := map[string]interface{}{
		"aps": map[string]interface{}{
			"alert":             description,
			"sound":             "default",
			"badge":             1,
			"content-available": 1,
		},
		"type":       types,
		"id":         id,
		"msgContent": description,
		"title":      title,
	}
	b, err := json.Marshal(out)
	if err != nil {
		cfg.LogErr("err:", err)

		return ""
	}
	return string(b)
}
