/**
 *
 * User: zhangbob
 * Date: 2019-05-28
 * Time: 19:12
 */
package control

import (
	"jiaotou.com/appPushSystem/config"
	"jiaotou.com/appPushSystem/model"
	app2 "jiaotou.com/appPushSystem/pushCore/app"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	args := []string{"welfareOpenPlatform", "--dir", "/Users/zhangbob/git/golang/jiaotou/appPushSystem", "--is_test", "true"}

	app := config.GetApp()
	config.SetAfter(func() {
		model.StartDb()
		go app2.StartPushCore()
	})
	app.Run(args)

	os.Exit(m.Run())
}

func TestSinglePush(t *testing.T) {
	go singlePush()

	time.Sleep(time.Second * 60)
}
