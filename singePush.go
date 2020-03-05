/**
 * 单条推送队列.
 *
 * User: zhangbob
 * Date: 2018-10-30
 * Time: 18:31
 */
package appPushSystem

import (
	"appPushSystem/app/android"
	"appPushSystem/app/iOS"
	"errors"
)

type AndroidParam struct {
	Key    string
	AppId  string
	Secret string
}

type IosParam struct {
	Pem         []byte
	Password    string
	PackageName string
	IsDev       bool
}

func NewSinglePushHandle(and AndroidParam, ios IosParam) (*singlePush, error) {
	m := &singlePush{messageList: make(chan messageParam, 500000)}
	m.androidInstance = android.NewAndroidPush(and.Key, and.AppId, and.Secret)
	err, iOSInstance := iOS.NewIosPush(ios.Pem, ios.Password, ios.PackageName, ios.IsDev)
	if err != nil {
		return nil, err
	}

	m.iOSInstance = iOSInstance

	go m.loop()

	return m, nil
}

type singlePush struct {
	messageList     chan messageParam
	androidInstance IMessageHandle
	iOSInstance     IMessageHandle
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
func (s *singlePush) AddSingleMessage(message IMessage, callback IHandleMessageCallback) {
	callback.Sending(message) //发送中
	s.messageList <- messageParam{Message: message, Callback: callback}
}

func (s singlePush) loop() {
	for v := range s.messageList {
		if instance := s.getSystem(v.Message.GetSystemType()); instance != nil {
			instance.PushSingle(v.Message, v.Callback)
		} else {
			v.Callback.Fail(v.Message, errors.New("手机类型有误"))
		}
	}
}

func (s *singlePush) getSystem(pushType int) IMessageHandle {
	if pushType == SystemAndroid {
		return s.androidInstance
	}

	if pushType == SystemIos {
		return s.iOSInstance
	}

	return nil
}

type messageParam struct {
	Message  IMessage
	Callback IHandleMessageCallback
}
