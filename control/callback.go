/**
 * 实现回调
 *
 * User: zhangbob
 * Date: 2020/3/4
 * Time: 9:49 上午
 */
package control

import (
	"sbjr.com/appPushSystem/config"
	"sbjr.com/appPushSystem/model"
	"sbjr.com/appPushSystem/pushCore"
)

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
