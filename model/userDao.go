/**
 * 用户信息.
 *
 * User: zhangbob
 * Date: 2018/5/9
 * Time: 下午2:10
 */
package model

import (
	"jiaotou.com/appPushSystem/pkg/cfg"
)

//获取所有用户
func GetAllUser(page, pageSize int) []UserInfo {
	var temp []UserInfo
	err := GetDb().Limit(pageSize, (page-1)*pageSize).Find(&temp)
	if err != nil {
		cfg.LogWarn("获取所有用户的ID失败: ", err.Error())
		return nil
	}

	return temp
}
