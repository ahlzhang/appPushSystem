/**
 * 单条推送队列.
 *
 * User: zhangbob
 * Date: 2018-10-30
 * Time: 18:31
 */
package app

import (
	"errors"
	"jiaotou.com/appPushSystem/pushCore"
	"jiaotou.com/appPushSystem/pushCore/app/android"
	"jiaotou.com/appPushSystem/pushCore/app/iOS"
)

var instance SinglePush

//启动单条推送队列
func StartPushCore() {
	instance = SinglePush{}
	instance.messageList = make(chan messageParam, 500000)
	instance.loop()
}

/**
 * 添加单条推送
 *
 * @param pushId 推送内容的索引值
 * @param token 推送手机的唯一标实符
 * @param content 推送内容
 * @param systemType 推送手机类型。1:android 2:iOS
 * @param okCallback 成功回调。(方法)
 * @param failCallback 失败回调。(方法)
 *
 **/
func AddSingleMessage(message pushCore.IMessage, callback pushCore.IHandleMessageCallback) {
	m := messageParam{Message: message, Callback: callback}

	//callback.Sending(message) //发送中
	instance.messageList <- m
}

type SinglePush struct {
	messageList chan messageParam
}

func (s SinglePush) loop() {
	for v := range s.messageList {
		if instance := getSystem(v.Message.GetSystemType()); instance != nil {
			instance.PushSingle(v.Message, v.Callback)
		} else {
			v.Callback.Fail(v.Message, errors.New("手机类型有误"))
		}
	}
}

type messageParam struct {
	Message  pushCore.IMessage
	Callback pushCore.IHandleMessageCallback
}

func getSystem(pushType int) pushCore.IMessageHandle {
	if pushType == pushCore.SystemAndroid {
		return android.GetAndroidHandleInstance()
	}

	if pushType == pushCore.SystemIos {
		return iOS.GetIosHandleInstance()
	}

	return nil
}
