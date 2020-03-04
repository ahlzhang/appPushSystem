/**
 * 单推
 *
 * User: zhangbob
 * Date: 2019-04-30
 * Time: 11:11
 */
package control

import (
	"sbjr.com/appPushSystem/pkg/cfg"
	"sbjr.com/appPushSystem/pushCore"
	"time"
)

//消息推送线程
//按照自己的需求进行调整时间
func StartPushMessageLoop() {
	cfg.LogInfo("推送处理队列启动...")

	instance := pushCore.GetCoreInstance(&HandleMessages{}, defaultCallback)

	t := time.NewTimer(time.Second * 10)
	for range t.C {
		instance.Push()
		t.Reset(time.Minute * 1)
	}
}
