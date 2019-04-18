package igetui

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
)

type IGeTui struct {
	Host         string
	AppKey       string
	MasterSecret string
}

func NewIGeTui(host, appkey, mastersecret string) *IGeTui {
	return &IGeTui{
		Host:         host,
		AppKey:       appkey,
		MasterSecret: mastersecret,
	}
}

func (iGeTui *IGeTui) connect() bool {
	sign := iGeTui.GetSign(iGeTui.AppKey, iGeTui.GetCurrentTime(), iGeTui.MasterSecret)
	params := map[string]interface{}{}
	params["action"] = "connect"
	params["appkey"] = iGeTui.AppKey
	params["timeStamp"] = iGeTui.GetCurrentTime()
	params["sign"] = sign

	rep := iGeTui.HttpPost(params)

	if "sucess" == rep["result"] {
		return true
	} else {
		fmt.Println("connect failed")
	}

	return false
}

func (iGeTui *IGeTui) PushMessageToSingle(message IGtSingleMessage, tartget Target) map[string]interface{} {
	params := map[string]interface{}{}
	params["action"] = "pushMessageToSingleAction"
	params["appkey"] = iGeTui.AppKey
	transparent := message.Data.GetTransparent()
	fmt.Println(transparent)
	byteArray, _ := proto.Marshal(transparent)
	params["clientData"] = base64.StdEncoding.EncodeToString(byteArray)
	params["transmissionContent"] = message.Data.GetTransmissionContent()
	params["isOffline"] = message.IsOffline
	params["offlineExpireTime"] = message.OfflineExpireTime
	params["appId"] = tartget.AppId
	params["clientId"] = tartget.ClientId
	params["type"] = 2
	params["pushType"] = message.Data.GetPushType()

	return iGeTui.HttpPostJson(params)

}

func (iGeTui *IGeTui) PushMessageToApp(message IGtAppMessage) map[string]interface{} {
	params := map[string]interface{}{}
	params["action"] = "pushMessageToAppAction"
	params["appkey"] = iGeTui.AppKey
	transparent := message.Data.GetTransparent()
	fmt.Println(transparent)
	byteArray, _ := proto.Marshal(transparent)
	params["clientData"] = base64.StdEncoding.EncodeToString(byteArray)
	params["transmissionContent"] = message.Data.GetTransmissionContent()
	params["isOffline"] = message.IsOffline
	params["offlineExpireTime"] = message.OfflineExpireTime
	params["appIdList"] = message.AppIdList
	params["phoneTypeList"] = message.PhoneTypeList
	params["provinceList"] = message.ProvinceList
	params["type"] = 2
	params["pushType"] = message.Data.GetPushType()

	return iGeTui.HttpPostJson(params)
}

func (iGeTui *IGeTui) PushMessageToList(contentId string, targets []Target) map[string]interface{} {
	params := map[string]interface{}{}
	params["action"] = "pushMessageToListAction"
	params["appkey"] = iGeTui.AppKey
	params["contentId"] = contentId

	targetList := []interface{}{}
	for _, target := range targets {
		appId := target.AppId
		clientId := target.ClientId
		targetTmp := map[string]string{"appId": appId, "clientId": clientId}
		targetList = append(targetList, targetTmp)
	}

	params["targetList"] = targetList
	params["type"] = 2

	return iGeTui.HttpPostJson(params)
}

func (iGeTui *IGeTui) GetContentId(message IGtListMessage) interface{} {
	params := map[string]interface{}{}
	params["action"] = "getContentIdAction"
	params["appkey"] = iGeTui.AppKey
	transparent := message.Data.GetTransparent()
	byteArray, _ := proto.Marshal(transparent)
	params["clientData"] = base64.StdEncoding.EncodeToString(byteArray)
	params["transmissionContent"] = message.Data.GetTransmissionContent()
	params["isOffline"] = message.IsOffline
	params["offlineExpireTime"] = message.OfflineExpireTime
	params["pushType"] = message.Data.GetPushType()
	ret := iGeTui.HttpPostJson(params)

	if ret["result"] == "ok" {
		return ret["contentId"]
	} else {
		return " "
	}
}

func (iGeTui *IGeTui) cancelContentId(contentId string) bool {
	params := map[string]interface{}{}
	params["action"] = "cancelContentIdAction"
	params["contentId"] = contentId

	ret := iGeTui.HttpPostJson(params)

	if ret["result"] == "ok" {
		return true
	} else {
		return false
	}
}

func (iGeTui *IGeTui) GetSign(appKey string, timeStamp int64, masterSecret string) string {
	rawValue := appKey + strconv.FormatInt(timeStamp, 10) + masterSecret
	h := md5.New()
	io.WriteString(h, rawValue)
	return hex.EncodeToString(h.Sum(nil))
}

func (iGeTui *IGeTui) GetCurrentTime() int64 {
	t := time.Now().Unix() * 1000
	return t
}

func (iGeTui *IGeTui) HttpPostJson(params map[string]interface{}) map[string]interface{} {
	ret := iGeTui.HttpPost(params)
	if ret["result"] == "sign_error" {
		iGeTui.connect()
		ret = iGeTui.HttpPost(params)
	}
	return ret
}

func (iGeTui *IGeTui) HttpPost(params map[string]interface{}) map[string]interface{} {
	data, _ := json.Marshal(params)
	fmt.Printf("%s\n", data)
	tryTime := 1
tryAgain:
	res, err := http.Post(iGeTui.Host, "application/json", strings.NewReader(string(data)))
	if err != nil {
		fmt.Println("第"+strconv.Itoa(tryTime)+"次", "请求失败")
		tryTime += 1
		if tryTime < 4 {
			goto tryAgain
		}
		return map[string]interface{}{"result": "post error"}
	}

	body, _ := ioutil.ReadAll(res.Body)
	var ret map[string]interface{}
	json.Unmarshal(body, &ret)
	return ret
}
