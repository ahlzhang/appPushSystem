/**
 * 单推
 *
 * User: zhangbob
 * Date: 2019-04-30
 * Time: 11:11
 */
package control

import (
	"encoding/json"
	"sbjr.com/appPushSystem/config"
	"sbjr.com/appPushSystem/model"
	"sbjr.com/appPushSystem/pkg/cfg"
	"sbjr.com/appPushSystem/pushCore"
	"time"
)

//消息推送线程
//按照自己的需求进行调整时间
func StartPushMessageLoop() {
	cfg.LogInfo("推送处理队列启动...")

	instance := pushCore.GetCoreInstance(&HandleMessages{}, defaultCallback)

	t := time.NewTimer(time.Second * 10)
	for range t.C {
		instance.Push()
		t.Reset(time.Minute * 1)
	}
}

type HandleMessages struct {
	MessageList []pushCore.IMessage
}

//实现获取消息接口
func (t *HandleMessages) GetMessage() []pushCore.IMessage {
	// TODO 获取需要推送的消息；按照自己的实际情况实现
	p := model.GetPushList()
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

		t.MessageList = append(t.MessageList, message{v})
	}

	return t.MessageList
}

//转换成推送使用的类型。需要自定义
type message struct {
	model.PushStruct
}

func (m message) ToMessage() (int64, string, string) {
	var content string
	switch m.GetSystemType() {
	case config.SystemIos:
		content = iosEncode(m.Id, int64(m.TypeId), m.Content, m.Title)
	case config.SystemAndroid:
		content = androidEncode(m.Id, int64(m.TypeId), m.Content, m.Title)
	}

	return m.Id, m.PushTo, content
}

func (m message) GetSystemType() int {
	return m.AppType
}

//传给android端的消息样式
func androidEncode(id int64, types int64, description string, title string) string {
	out := map[string]interface{}{
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

//传给iOS端的消息样式
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