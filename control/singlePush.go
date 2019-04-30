/**
 * 单推
 *
 * User: zhangbob
 * Date: 2019-04-30
 * Time: 11:11
 */
package control

import (
	"jiaotou.com/appPushSystem/config"
	"jiaotou.com/appPushSystem/model"
	"jiaotou.com/appPushSystem/pkg/cfg"
	"jiaotou.com/appPushSystem/pushCore"
	"time"
)

//消息推送线程
func StartPushMessageLoop() {
	cfg.LogInfo("推送处理队列启动...")

	t := time.NewTimer(time.Second * 10)
	for range t.C {
		singlePush()
		t.Reset(time.Minute * 1)
	}
}

func singlePush() {
	PushMessage(func() ([]pushCore.IMessage, pushCore.IHandleMessageCallback) {
		p := model.GetPushList()

		var msgList []pushCore.IMessage

		android := config.SystemAndroid
		ios := config.SystemIos
		for _, v := range p {
			if v.AppType != pushCore.SystemAndroid && v.AppType != pushCore.SystemIos {
				continue
			}

			if v.AppType == android {
				if v.AppType != pushCore.SystemAndroid {
					v.AppType = pushCore.SystemAndroid
				}
			}

			if v.AppType == ios {
				if v.AppType != pushCore.SystemIos {
					v.AppType = pushCore.SystemIos
				}
			}

			if v.PushTo == "" {
				continue
			}

			msgList = append(msgList, message{v})
		}

		return msgList, nil
	})
}
