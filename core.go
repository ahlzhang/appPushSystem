/**
 * 核心类
 *
 * User: zhangbob
 * Date: 2020/3/3
 * Time: 5:17 下午
 */
package appPushSystem

//创建实例，项目只能有一个此实例
func GetCoreInstance(message GetMessage, callback IHandleMessageCallback, and AndroidParam, ios IosParam) (*core, error) {
	instance := &core{
		GetMessage:             message,
		IHandleMessageCallback: callback,
	}

	in, err := NewSinglePushHandle(and, ios)
	if err != nil {
		return nil, err
	}

	instance.pushHandle = in

	return instance, nil
}

// 获取消息
type GetMessage interface {
	GetMessage() []IMessage
}

type core struct {
	GetMessage
	IHandleMessageCallback
	pushHandle *singlePush
}

func (t *core) Push() {
	list := t.GetMessage.GetMessage()

	if len(list) == 0 {
		return
	}

	for _, v := range list {
		t.pushHandle.AddSingleMessage(v, t.IHandleMessageCallback)
	}
}
