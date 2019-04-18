/**
 * grpc 入口.
 *
 * User: zhangbob
 * Date: 2018/5/8
 * Time: 下午5:53
 */
package api

import (
	"context"
	"encoding/json"
	"jiaotou.com/appPushSystem/config"
	"jiaotou.com/appPushSystem/control"
	"jiaotou.com/appPushSystem/grpc"
	"jiaotou.com/appPushSystem/model"
	"jiaotou.com/appPushSystem/pkg/cfg"
)

const (
	addPushMessage = 100 //添加系统推送信息
)

var funcMap map[int32]func(string) (int, string, string)

func StartFunLoad() {
	funcMap = make(map[int32]func(string) (int, string, string))
	funcMap[addPushMessage] = addPushMessageApi
}

type PushRequestIml struct {
}

func (PushRequestIml) PushServiceRequest(cxt context.Context, in *grpc.RequestParam) (*grpc.RequestResult, error) {
	result := new(grpc.RequestResult)

	if f, ok := funcMap[in.Code]; ok {
		cfg.LogInfof("====>grpc调用入参:code:%d,参数:%s", in.Code, in.Param)

		code, msg, re := f(in.Param)
		if code != 0 {
			result.Success = 0
			result.ErrorCode = int32(code)
			result.ErrorMessage = msg
		} else {
			result.Success = 1
			result.Result = re
		}
	} else {
		result.Success = 0
		result.ErrorCode = 2
		result.ErrorMessage = "参数错误。code无效"
	}

	cfg.LogInfof("====>grpc调用返回结果:success:%d,code:%d,message:%s,result:%s", result.Success, result.ErrorCode, result.ErrorMessage, result.Result)

	return result, nil
}

/**
 * @api grpc /100 推送系统信息
 * @apiGroup 推送系统
 * @apiParam title string 标题
 * @apiParam url string 跳转地址
 * @apiParam img string 图片地址
 * @apiParam typeId int 信息类型
 * @apiParam content string 信息内容
 * @apiParam channel int 消息发送渠道 1app 2运营平台
 * @apiParam pushRange int 消息范围 1所有人 2指定人群
 * @apiParam userIdList []int 需要推送的用户ID集合
 *
 * @apiSuccess 200 json ok
 * @apiExample json
 * 例子
 **/
func addPushMessageApi(param string) (int, string, string) {
	p := struct {
		Title      string  `json:"title"`
		Img        string  `json:"img"`
		Url        string  `json:"url"`
		TypeId     int     `json:"typeId"`
		Content    string  `json:"content"`
		Channel    int     `json:"channel"`
		PushRange  int     `json:"pushRange"`
		UserIdList []int64 `json:"userIdList"`
	}{}

	err := json.Unmarshal([]byte(param), &p)
	if err != nil {
		return 30, "json解析失败", ""
	}

	if p.PushRange != config.PushRangeCustom && p.PushRange != config.PushRangeAll {
		return 31, "pushRange错误", ""
	}

	if p.PushRange == config.PushRangeCustom && len(p.UserIdList) == 0 {
		return 31, "推送范围错误", ""
	}

	msg, _ := model.GetMessageType(p.TypeId)
	if msg != "" {
		return 31, msg, ""
	}

	go control.AddMessageControl(p.Title, p.Img, p.Url, p.Content, p.Channel, p.PushRange, p.TypeId, p.UserIdList)
	return 0, "", ""
}
