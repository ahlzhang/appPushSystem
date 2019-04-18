/**
 * 推送信息核心.
 *
 * User: zhangbob
 * Date: 2018/5/9
 * Time: 上午9:21
 */
package control

import (
	apns "github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"jiaotou.com/appPushSystem/config"
	"jiaotou.com/appPushSystem/model"
	"jiaotou.com/appPushSystem/pkg/cfg"
	"time"
)

type ListSendMsg struct {
	PushId  int64
	Token   string
	Content string
}

const allPushTime = time.Minute * 5

func initProperty() {
	androidPushList = make(chan ListSendMsg, 100000)
	iosPushList = make(chan ListSendMsg, 100000)

	cert, err := certificate.FromPemFile(config.AppConfPath+"/conf/"+config.Conf.PushInfo.Ios.Pem, config.Conf.PushInfo.Ios.Password)
	if err != nil {
		cfg.LogFatal("ios证书不正确:", err)
	}

	if config.Conf.PushInfo.Ios.Model == "PRO" {
		cfg.LogInfo("PRO环境:", err)
		apnsClient = apns.NewClient(cert).Production()
	} else {
		cfg.LogInfo("DEV环境:", err)
		apnsClient = apns.NewClient(cert).Development()
	}

	go sendLoopAnd()
	go sendLoopIos()
}

//开启推送服务线程
func StartTask() {
	initProperty()

	cfg.LogInfo("推送服务开始开启:")

	sec_2 := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-sec_2.C:
			loop()
			sec_2.Reset(time.Second * 1)
		}
	}
}

func loop() {
	customPush()
	allPush()
}

//单条与指定人群推送
func customPush() {
	list := model.GetPushList()

	customPush := make(map[int64][]struct {
		Cid     string
		Id      int64
		AppType int
	})
	customMap := make(map[int64]model.PushStruct)

	for _, v := range list {
		if v.PushTo == "" {
			continue
		}

		customPush[v.MessageId] = append(customPush[v.MessageId], struct {
			Cid     string
			Id      int64
			AppType int
		}{Cid: v.PushTo, Id: v.Id, AppType: v.AppType})
		customMap[v.MessageId] = v
	}

	for k, v := range customPush {
		message := customMap[k]
		if len(v) > 1 {
			var cidList []string
			var idList []int64
			for _, v1 := range v {
				if v1.AppType == config.SystemIos {
					m := message
					m.PushTo = v1.Cid
					m.Id = v1.Id
					sendPush(m)
				} else {
					cidList = append(cidList, v1.Cid)
					idList = append(idList, v1.Id)
				}
			}
			PushList_Android(cidList, idList, androidEncode(message.MessageId, int64(message.TypeId), message.Content, message.Title))
			model.UpdateCenterByMessageId(message.MessageId, config.PushStateAlready)
		} else {
			sendPush(message)
			model.UpdateCenterByMessageId(message.MessageId, config.PushStateAlready)
		}
	}
}

func allPush() {
	list := model.GetAllPushList()
	if len(list) == 0 {
		return
	}

	for _, v := range list {
		var appIdList []string
		var idList []int64
		if time.Unix(v.CreateTime, 0).Add(allPushTime).Before(time.Now()) {

			for i := 1; i == i; i++ {
				userList := model.GetAllUser(i, 5000)
				if len(userList) == 0 {
					break
				}

				var messageUser []model.MessageUser
				for _, u := range userList {
					messageUser = append(messageUser, model.MessageUser{
						UserId:     u.UserId,
						CreateTime: time.Now().Unix(),
						PushState:  config.PushStateWaiting,
						MessageId:  v.MessageId,
						IsDel:      config.DelStateNot,
						ReadState:  config.ReadStateNotRead,
					})
				}

				model.InsertMessageUser(messageUser)

				iosList := model.GetPushListByMessageId(config.SystemIos, v.MessageId)
				if len(iosList) > 0 {
					for _, v1 := range iosList {
						sendPush(v1)
					}
				}

				andList := model.GetPushListByMessageId(config.SystemAndroid, v.MessageId)
				for _, v := range andList {
					appIdList = append(appIdList, v.PushTo)
					idList = append(idList, v.Id)
				}
			}

			PushList_Android(appIdList, idList, androidEncode(v.MessageId, int64(v.TypeId), v.Content, v.Title))
			model.UpdateCenterByMessageId(v.MessageId, config.PushStateAlready)
		}
	}

}

func sendPush(obj model.PushStruct) {
	model.UpdatePushBySingle(obj.Id, config.PushStateSending)
	if obj.AppType == config.SystemAndroid {
		content := androidEncode(obj.MessageId, int64(obj.TypeId), obj.Content, obj.Title)
		Push_Android(obj.Id, obj.PushTo, content)
		msg := ListSendMsg{
			PushId:  obj.Id,
			Content: content,
			Token:   obj.PushTo,
		}
		androidPushList <- msg

	} else if obj.AppType == config.SystemIos {
		content := iosEncode(obj.MessageId, int64(obj.TypeId), obj.Content, obj.Title)
		msg := ListSendMsg{
			PushId:  obj.Id,
			Content: content,
			Token:   obj.PushTo,
		}
		iosPushList <- msg
	}
}
