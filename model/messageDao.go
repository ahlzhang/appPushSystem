package model

import (
	"jiaotou.com/appPushSystem/pkg/cfg"
	"time"
	"fmt"
	"jiaotou.com/appPushSystem/config"
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

func UpdatePushByGroup(pushIdList []int64, status int) bool {
	_, err := GetDb().Cols("push_state", "push_time").In("id", pushIdList).
		Update(MessageUser{PushTime: time.Now().Unix(), PushState: status})

	if err != nil {
		cfg.LogErr("err:", err)

		return false
	}
	return true
}

func UpdatePushByMessageId(messageId int64, status int) bool {
	_, err := GetDb().Cols("push_state", "push_time").Where("message_id = ?", messageId).
		Update(MessageUser{PushTime: time.Now().Unix(), PushState: status})

	if err != nil {
		cfg.LogErr("err:", err)

		return false
	}
	return true
}

func UpdateCenterByMessageId(messageId int64, status int) bool {
	_, err := GetDb().Cols("status").Where("message_id = ?", messageId).
		Update(MessagesCenter{Status: status})

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
			c.push_range,
			c.create_time,
			u.cid,
			u.operating_system_type as app_type
		FROM
			t_app_message_center c
			LEFT JOIN t_app_message_user_state s ON c.message_id = s.message_id
			LEFT JOIN t_app_user u ON s.user_id = u.user_id
		WHERE s.push_state = %d AND c.push_range <> %d`, config.PushStateWaiting, config.PushRangeAll)

	err := GetDb().SQL(sql).Find(&pushList)
	if err != nil {
		cfg.LogErr("推送列表查询失败:", err)
	}

	return
}

func GetAllPushList() []MessagesCenter {
	var r []MessagesCenter
	err := GetDb().Where("push_range = ?", config.PushRangeAll).And("status = ?", config.PushStateWaiting).Find(&r)
	if err != nil {
		cfg.LogErr("获取群体推送出错.", err.Error())
		return nil
	}

	return r
}

//根据messageId获取未推送的信息
func GetPushListByMessageId(systemType int, messageId int64) (pushList []PushStruct) {
	sql := fmt.Sprintf(
		`SELECT
		  s.id,
			s.user_id,
			c.message_id,
			c.title,
			c.type_id,
			c.content,
			c.push_range,
			c.create_time,
			u.cid,
			u.operating_system_type as app_type
		FROM
			t_app_message_center c
			LEFT JOIN t_app_message_user_state s ON c.message_id = s.message_id
			LEFT JOIN t_app_user u ON s.user_id = u.user_id
		WHERE s.push_state = %d AND c.message_id = %d AND u.operating_system_type = %d`,
		config.PushStateWaiting, messageId, systemType)

	err := GetDb().SQL(sql).Find(&pushList)
	if err != nil {
		cfg.LogErr("推送列表查询失败:", err)
	}

	return
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

//func GetUserMessageBySystem(userId int, page, pageSize int) (string, []MessageView) {
//	var result []MessageView
//
//	err := GetDb().Where("is_del = 0 AND type_id = 0 AND user_id = ?", userId).Limit(pageSize, (page-1)*pageSize).Find(&result)
//	if err != nil {
//		cfg.LogErr("err:", err)
//
//		return "服务器异常请稍后再试", nil
//	}
//	return "", result
//}
//
//func GetUserMessageByActive(userId int, page, pageSize int) (string, []MessageView) {
//	var result []MessageView
//
//	err := GetDb().Where("is_del = 0 AND type_id = 1 AND user_id = ?", userId).Limit(pageSize, (page-1)*pageSize).Find(&result)
//	if err != nil {
//		cfg.LogErr("err:", err)
//
//		return "服务器异常请稍后再试", nil
//	}
//	return "", result
//}
//
//func GetUserHotMessage(userId int) (string, []MessageView) {
//	var result []MessageView
//
//	err := GetDb().Where("is_del = 0 AND type_id = 1 AND read_state = 0 AND user_id = ?", userId).Find(&result)
//	if err != nil {
//		cfg.LogErr("err:", err)
//
//		return "服务器异常请稍后再试", result
//	}
//
//	return "", result
//}
//
//func GetUserNotReadMessage(userId int) (string, int) {
//	var result []MessageView
//	err := GetDb().Where("read_state = 0").And("user_id = ?", userId).Find(&result)
//
//	if err != nil {
//		cfg.LogErr("err:", err)
//
//		return "服务器异常请稍后再试", 0
//	}
//
//	return "", len(result)
//}
//
////更新已读状态。如果idList为空则是全部更新。
//func UpdateMessageReadState(userId int, idList []int64) bool {
//	var p MessageUser
//	p.ReadState = 1
//	p.ReadTime = time.Now().Unix()
//
//	var err error
//	if len(idList) == 0 {
//		_, err = GetDb().Cols("read_time", "read_state").Where("user_id = ?", userId).Update(p)
//	} else {
//		_, err = GetDb().Cols("read_time", "read_state").Where("user_id = ?", userId).In("id", idList).Update(p)
//	}
//
//	if err != nil {
//		cfg.LogErr("err:", err)
//
//		return false
//	}
//
//	return true
//}
//
////更新删除状态。如果idList为空则是全部更新。
//func UpdateMessageDelState(userId int, idList []int64) bool {
//	var p MessageUser
//	p.IsDel = 1
//	p.ReadState = 1
//	p.ReadTime = time.Now().Unix()
//
//	var err error
//	if len(idList) == 0 {
//		_, err = GetDb().Cols("is_del", "read_state", "read_time").Where("user_id = ?", userId).Update(p)
//	} else {
//		_, err = GetDb().Cols("is_del", "read_state", "read_time").Where("user_id = ?", userId).In("id", idList).Update(p)
//	}
//
//	if err != nil {
//		cfg.LogErr("err:", err)
//
//		return false
//	}
//
//	return true
//}
