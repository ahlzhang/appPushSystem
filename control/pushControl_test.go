/** 
 * 测试
 *
 * User: zhangbob 
 * Date: 2018/6/26 
 * Time: 下午5:30 
 */
package control

import (
	"testing"
	"jiaotou.com/appPushSystem/model"
	"fmt"
		"os"
		)

func TestMain(m *testing.M) {
	//time.Sleep(time.Second * 5)
	os.Exit(m.Run())
}

func TestPushIos(t *testing.T) {
	//sql := fmt.Sprintf(
	//	`SELECT
	//	  s.id,
	//		s.user_id,
	//		c.message_id,
	//		c.title,
	//		c.type_id,
	//		c.content,
	//		c.push_range,
	//		c.create_time,
	//		u.cid,
	//		u.operating_system_type as app_type
	//	FROM
	//		t_app_message_center c
	//		LEFT JOIN t_app_message_user_state s ON c.message_id = s.message_id
	//		LEFT JOIN t_app_user u ON s.user_id = u.user_id
	//	WHERE s.id = %d AND u.operating_system_type = 2 `, 20697)

	var pushList model.PushStruct
	//isExist, err := model.MysqlDefault.SQL(sql).Get(&pushList)
	//if err != nil {
	//	t.Error("========>出错:", err.Error())
	//	return
	//}
	//if !isExist {
	//	t.Error("========>不存在")
	//	return
	//}
	pushList.MessageId = 1
	pushList.TypeId = 1
	pushList.Content = "小李子，测试推送"
	pushList.Title = "小李子"
	pushList.Id = 2
	pushList.PushTo = "f6ffbf6bcaebc3d160c13fd8316455a3b1d4374ed3e4ae2e6f2fd7c46f225e60"
	content := iosEncode(pushList.MessageId, int64(pushList.TypeId), pushList.Content, pushList.Title)

	PushIos(pushList.Id, pushList.PushTo, content)
}

func TestPushAndroid(t *testing.T) {
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
		WHERE s.id = %d AND u.operating_system_type = 1 `, 20545)

	var pushList model.PushStruct
	isExist, err := model.MysqlDefault.SQL(sql).Get(&pushList)
	if err != nil {
		t.Error("========>出错:", err.Error())
		return
	}
	if !isExist {
		t.Error("========>不存在")
		return
	}
	content := androidEncode(pushList.MessageId, int64(pushList.TypeId), pushList.Content, pushList.Title)
	Push_Android(pushList.Id, pushList.PushTo, content)
}

func TestPushAll(t *testing.T) {
	allPush()
}

func TestPushCustom(t *testing.T) {
	customPush()
}
