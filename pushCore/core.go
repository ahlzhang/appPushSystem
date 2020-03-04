/**
 * 核心类
 *
 * User: zhangbob
 * Date: 2020/3/3
 * Time: 5:17 下午
 */
package pushCore

import (
	"sbjr.com/appPushSystem/pushCore/app"
	"sync"
)

var instance *core
var once sync.Once

func GetCoreInstance(message GetMessage, callback IHandleMessageCallback) *core {
	once.Do(func() {
		instance = &core{
			GetMessage:             message,
			IHandleMessageCallback: callback,
		}

		go app.StartPushCore()
	})

	return instance
}

// 获取消息
type GetMessage interface {
	GetMessage() []IMessage
}

type core struct {
	GetMessage
	IHandleMessageCallback
}

func (t *core) Push() {
	list := t.GetMessage.GetMessage()

	if len(list) == 0 {
		return
	}

	for _, v := range list {
		app.AddSingleMessage(v, t.IHandleMessageCallback)
	}
}
