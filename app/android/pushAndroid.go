/**
 * android推送.
 *
 * User: zhangbob
 * Date: 2018/5/8
 * Time: 下午6:15
 */
package android

import (
	"errors"
	"fmt"
	"github.com/ahlzhang/appPushSystem"
	"github.com/ahlzhang/appPushSystem/app/android/igetui"
	"github.com/ahlzhang/appPushSystem/app/android/igetui/template"
)

const igtHost = "http://sdk.open.api.igexin.com/apiex.htm"

func NewAndroidPush(appKey, appId, secret string) *messageHandle {
	return &messageHandle{
		key:    appKey,
		appId:  appId,
		secret: secret,
	}
}

type messageHandle struct {
	key    string
	appId  string
	secret string
}

func (t messageHandle) PushSingle(msg appPushSystem.IMessage, callBack appPushSystem.IHandleMessageCallback) {
	_, clientId, content := msg.ToMessage()
	if clientId == "" {
		callBack.Fail(msg, errors.New("tokenId为空"))
		return
	}

	pushs := igetui.NewIGeTui(igtHost, t.key, t.secret)
	data := template.NewTransmissionTemplate(t.appId, t.key, 2, content)
	message := igetui.NewIGtSingleMessage(true, 3600, data)
	target := igetui.NewTarget(t.appId, clientId)
	result := pushs.PushMessageToSingle(*message, *target)

	if v, ok := result["result"]; ok && v == "ok" {
		callBack.Success(msg)
	} else {
		callBack.Fail(msg, errors.New(fmt.Sprintf("安卓推送失败:%s", result)))
	}
}

//群推。暂时没用。因为如果推送后，messageId无法关联。
func (t messageHandle) PushList(cIdList []string, idList []int64, content string) error {
	var targetList []igetui.Target
	for _, v := range cIdList {
		temp := igetui.NewTarget(t.appId, v)
		targetList = append(targetList, *temp)
	}

	pushs := igetui.NewIGeTui(igtHost, t.key, t.secret)
	data := template.NewTransmissionTemplate(t.appId, t.key, 2, content)
	message := igetui.NewIGtListMessage(true, 3600, data)

	contentId := pushs.GetContentId(*message)
	if contentId == " " {
		return errors.New("获取contentId失败")
	}

	result := pushs.PushMessageToList(contentId.(string), targetList)
	if v, ok := result["result"]; ok && v == "ok" {
		return nil
	} else {
		return errors.New(fmt.Sprintf("安卓推送失败,返回值:%s", result))
	}
}
