package model

import (
	"fmt"
	"sbjr.com/appPushSystem/config"
	"sbjr.com/appPushSystem/pkg/cfg"
	"time"
)

func UpdatePushBySingle(pushId int64, status int) bool {
	_, err := GetDb().Cols("push_state", "push_time").Where("id = ?", pushId).
		Update(MessageUser{PushTime: time.Now().Unix(), PushState: status})
	if err != nil {
		cfg.LogErr("err:", err)

		return false
	}
	return true
}

//指定人群推送
func GetPushList() (pushList []PushStruct) {
	sql := fmt.Sprintf(
		`SELECT
		  s.id,
			s.user_id,
			c.message_id,
			c.title,
			c.type_id,
			c.content,
			c.create_time,
			u.cid,
			u.operating_system_type as app_type
		FROM
			t_app_message_center c
			LEFT JOIN t_app_message_user_state s ON c.message_id = s.message_id
			LEFT JOIN t_app_user u ON s.user_id = u.user_id
		WHERE s.push_state = %d and u.cid <> ''`, config.PushStateWaiting)

	err := GetDb().SQL(sql).Find(&pushList)
	if err != nil {
		cfg.LogErr("推送列表查询失败:", err)
	}

	return
}

//获取消息类型
func GetMessageType(typeId int) (string, MessageType) {
	var result MessageType

	exist, err := GetDb().Where("type_id = ?", typeId).Get(&result)
	if err != nil {
		cfg.LogErr("err:", err)
		return "操作失败,请稍后再试", result
	}

	if !exist {
		return "无找到相关的消息类型", result
	}

	return "", result
}

// 添加消息
func InsertMessage(param MessagesCenter) (MessagesCenter, string) {
	_, err := GetDb().Insert(&param)
	if err != nil {
		cfg.LogErr("err:", err)

		return param, "消息插入失败，请稍后再试"
	}

	return param, ""
}

// 添加对应的消息关联信息
func InsertMessageUser(param []MessageUser) (string) {
	_, err := GetDb().Insert(&param)
	if err != nil {
		cfg.LogErr("err:", err)

		return "用户与消息关联插入失败，请稍后再试"
	}

	return ""
}

func UpdateState(id []int64, state int) {
	_, err := GetDb().Cols("push_state").In("id", id).Update(MessageUser{PushState: state})
	if err != nil {
		cfg.LogWarn("批量更新商品详情失败:", err.Error())
	}
}
