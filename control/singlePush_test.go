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
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	args := []string{"welfareOpenPlatform", "--dir", "/Users/zhangbob/git/golang/welfareOpenPlatform", "--is_test", "true"}

	app := config.GetApp()
	config.SetAfter(func() {
		model.StartDb()
	})
	app.Run(args)

	os.Exit(m.Run())
}

func TestSinglePush(t *testing.T) {
	singlePush()
}
