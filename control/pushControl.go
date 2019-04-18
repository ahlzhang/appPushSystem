/** 
 * 推送api.
 *
 * User: zhangbob 
 * Date: 2018/5/22 
 * Time: 下午7:22 
 */
package control

import (
	"jiaotou.com/appPushSystem/model"
	"time"
	"jiaotou.com/appPushSystem/config"
)

func AddMessageControl(title, img, url, content string, channel, pushRange, typeId int, useridList []int64) {
	var contentType int
	if url == "" {
		contentType = config.ContentTypeNotJump
	} else {
		contentType = config.ContentTypeJumpOutside
	}

	m, msg := model.InsertMessage(model.MessagesCenter{
		PushRange:   pushRange,
		Content:     content,
		Url:         url,
		Img:         img,
		Title:       title,
		TypeId:      typeId,
		CreateTime:  time.Now().Unix(),
		ContentType: contentType,
		Channel:     channel,
	})

	if msg != "" {
		return
	}

	if pushRange == config.PushRangeAll {
		return
	}

	if len(useridList) == 0 {
		return
	}

	if len(useridList) < 5000 {
		messageUser := createMessageUser(useridList, m.MessageId)
		model.InsertMessageUser(messageUser)
	} else {
		length := len(useridList)
		for i := 1; i <= length/5000+1; i++ {
			start := (i - 1) * 5000
			end := start + 5000

			if start >= length {
				break
			}

			if end > length {
				end = length
			}

			model.InsertMessageUser(createMessageUser(useridList[start:end], m.MessageId))
		}
	}

	return
}

func createMessageUser(userIdList []int64, messageId int64) []model.MessageUser {
	var messageUser []model.MessageUser
	for _, v := range userIdList {
		messageUser = append(messageUser, model.MessageUser{
			UserId:     v,
			CreateTime: time.Now().Unix(),
			PushState:  config.PushStateWaiting,
			MessageId:  messageId,
			IsDel:      config.DelStateNot,
			ReadState:  config.ReadStateNotRead,
		})
	}

	return messageUser
}
