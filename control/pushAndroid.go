/**
 * android推送.
 *
 * User: zhangbob
 * Date: 2018/5/8
 * Time: 下午6:15
 */
package control

import (
	"encoding/json"
	"errors"
	"jiaotou.com/appPushSystem/config"
	"jiaotou.com/appPushSystem/control/igetui"
	"jiaotou.com/appPushSystem/control/igetui/template"
	"jiaotou.com/appPushSystem/model"
	"jiaotou.com/appPushSystem/pkg/cfg"
)

const igtHost = "http://sdk.open.api.igexin.com/apiex.htm"

var androidPushList = make(chan ListSendMsg)

func sendLoopAnd() {
	for {
		if v, ok := <-androidPushList; ok {
			Push_Android(v.PushId, v.Token, v.Content)
		}
	}
}

// 推送(单推)
func Push_Android(pushId int64, token, content string) error {
	pushs := igetui.NewIGeTui(igtHost, config.Conf.PushInfo.Android.AppKey, config.Conf.PushInfo.Android.MasterSecret)
	data := template.NewTransmissionTemplate(config.Conf.PushInfo.Android.AppId, config.Conf.PushInfo.Android.AppKey, 2, content)
	message := igetui.NewIGtSingleMessage(true, 3600, data)
	target := igetui.NewTarget(config.Conf.PushInfo.Android.AppId, token)
	result := pushs.PushMessageToSingle(*message, *target)

	if v, ok := result["result"]; ok && v == "ok" {
		model.UpdatePushBySingle(pushId, config.PushStateAlready)
	} else {
		cfg.LogErr("push error,  result:", result)
		model.UpdatePushBySingle(pushId, config.PushStateFail)
	}

	return nil
}

// 群发(同样的消息，同样的应用渠道)
func PushList_Android(token []string, idList []int64, content string) error {
	var targetList []igetui.Target
	for _, v := range token {
		temp := igetui.NewTarget(config.Conf.PushInfo.Android.AppId, v)
		targetList = append(targetList, *temp)
	}

	pushs := igetui.NewIGeTui(igtHost, config.Conf.PushInfo.Android.AppKey, config.Conf.PushInfo.Android.MasterSecret)
	data := template.NewTransmissionTemplate(config.Conf.PushInfo.Android.AppId, config.Conf.PushInfo.Android.AppKey, 2, content)
	message := igetui.NewIGtListMessage(true, 3600, data)

	contentId := pushs.GetContentId(*message)
	if contentId == " " {
		return errors.New("获取contentId失败")
	}

	model.UpdatePushByGroup(idList, config.PushStateSending)

	result := pushs.PushMessageToList(contentId.(string), targetList)
	if v, ok := result["result"]; ok && v == "ok" {
		model.UpdatePushByGroup(idList, config.PushStateAlready)
	} else {
		cfg.LogErr("push error,  result:", result)
		model.UpdatePushByGroup(idList, config.PushStateFail)
	}

	return nil
}

// 全部用户发布信息
func PushAll_Android(messageId int64, content string, token []string) error {
	model.UpdatePushByMessageId(messageId, config.PushStateSending)
	data := template.NewTransmissionTemplate(config.Conf.PushInfo.Android.AppId, config.Conf.PushInfo.Android.AppKey, 2, content)
	message := igetui.NewIGtAppMessage(true, 3600, data)
	message.AppIdList = token
	pushs := igetui.NewIGeTui(igtHost, config.Conf.PushInfo.Android.AppKey, config.Conf.PushInfo.Android.MasterSecret)
	result := pushs.PushMessageToApp(*message)
	if v, ok := result["result"]; ok && v == "post error" {
		model.UpdatePushByMessageId(messageId, config.PushStateFail)
	} else {
		model.UpdatePushByMessageId(messageId, config.PushStateAlready)
	}

	return nil
}

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
