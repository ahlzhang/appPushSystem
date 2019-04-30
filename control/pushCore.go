/**
 * 推送业务核心
 *
 * User: zhangbob
 * Date: 2019-04-30
 * Time: 10:50
 */
package control

import (
	"encoding/json"
	"jiaotou.com/appPushSystem/config"
	"jiaotou.com/appPushSystem/model"
	"jiaotou.com/appPushSystem/pkg/cfg"
	"jiaotou.com/appPushSystem/pushCore"
	"jiaotou.com/appPushSystem/pushCore/app"
)

func PushMessage(f func() ([]pushCore.IMessage, pushCore.IHandleMessageCallback)) {
	meg, callback := f()
	pushCore := HandleMessages{MessageList: meg, callback: callback}
	pushCore.Push()
}

type HandleMessages struct {
	callback    pushCore.IHandleMessageCallback
	MessageList []pushCore.IMessage
}

//单条消息处理。
func (h *HandleMessages) Push() {
	if len(h.MessageList) == 0 {
		return
	}

	if h.callback == nil {
		h.callback = defaultCallback
	}

	for _, v := range h.MessageList {
		app.AddSingleMessage(v, h.callback)
	}
}

//设置消息处理回调
func (h *HandleMessages) SetHandleCallback(callback pushCore.IHandleMessageCallback) {
	h.callback = callback
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

	return m.MessageId, m.PushTo, content
}

func (m message) GetSystemType() int {
	return m.AppType
}

var defaultCallback DefaultCallback

//默认消息处理回调.需要自己定义。默认为单条推送
type DefaultCallback struct {
}

func (DefaultCallback) Sending(message pushCore.IMessage) {
	id, _, _ := message.ToMessage()
	model.UpdatePushBySingle(id, config.PushStateSending)
}

func (DefaultCallback) Success(message pushCore.IMessage) {
	id, _, _ := message.ToMessage()
	model.UpdatePushBySingle(id, config.PushStateAlready)
}

func (DefaultCallback) Fail(message pushCore.IMessage, err error) {
	id, _, _ := message.ToMessage()
	model.UpdatePushBySingle(id, config.PushStateFail)
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
