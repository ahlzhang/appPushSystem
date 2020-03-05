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
	"sbjr.com/appPushSystem/grpc"
	"sbjr.com/appPushSystem/pkg/cfg"
)

var funcMap map[int32]func(string) (int, string, string)

func StartFunLoad() {
	funcMap = make(map[int32]func(string) (int, string, string))
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